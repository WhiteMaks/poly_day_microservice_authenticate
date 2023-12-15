package microservice

import (
	"context"
	"database/sql"
	"errors"
	"github.com/WhiteMaks/poly_day_bus_schema/schema"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/utils"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,max=25"`
	Password string `json:"password" binding:"required,max=72"`
}

type CreateUserResponse struct {
	Id int64 `json:"id"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type CheckUserResponse struct {
	Authorized bool `json:"authorized"`
}

var ErrAuthorizationDataNotValid = errors.New("authorization data is not correct")

const (
	ROLE_DEVELOPER = "ROLE_DEVELOPER"
)

func (ms *Microservice) checkUser(ctx *gin.Context) {
	response := CheckUserResponse{
		Authorized: false,
	}

	token := ctx.GetHeader("Authorization")
	if len(token) == 0 {
		ctx.JSON(http.StatusOK, response)
		return
	}

	_, err := ms.tokenCreator.Verify(token)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Authorized = true

	ctx.JSON(http.StatusOK, response)
}

func (ms *Microservice) loginUser(ctx *gin.Context) {
	var requestBody LoginUserRequest

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dbUser, err := ms.store.GetUserByUsername(ctx, requestBody.Username)

	userExist, err := checkExistDBRowError(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !userExist {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrAuthorizationDataNotValid))
		return
	}

	err = utils.PasswordEquality(requestBody.Password, dbUser.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrAuthorizationDataNotValid))
		return
	}

	token, err := ms.tokenCreator.Create(dbUser.Username, time.Hour*24*7)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := LoginUserResponse{
		Token: token,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ms *Microservice) createUser(ctx *gin.Context) {
	var requestBody CreateUserRequest

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = ms.store.GetUserByUsername(ctx, requestBody.Username)

	userExist, err := checkExistDBRowError(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if userExist {
		errorMessage := errors.New("user with the same name already exists")
		ctx.JSON(http.StatusConflict, errorResponse(errorMessage))
		return
	}

	dbRole, err := ms.store.GetRoleByName(ctx, ROLE_DEVELOPER)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(requestBody.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sqlUser := database.CreateUserParams{
		Username: requestBody.Username,
		Password: hashedPassword,
		RoleID:   dbRole.ID,
	}

	dbUser, err := ms.store.CreateUser(ctx, sqlUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := CreateUserResponse{
		Id: dbUser.ID,
	}

	message, err := schema.GenerateMessageForUserCreated10(dbUser.ID, dbRole.Name)
	if err != nil {
		log.Printf("error generate event user-created 1.0: %v", err)
	}

	err = ms.bus.SendMessageInExchange(schema.PathEventAuthenticateUserCreated10, message)
	if err != nil {
		log.Printf("error send event user-created 1.0: %v", err)
	}

	ctx.JSON(http.StatusCreated, response)
}

func checkExistDBRowError(err error) (bool, error) {
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, err
}

func (ms *Microservice) identifyUser10(message schema.Message) {
	responseMessage := schema.GenerateResponseMessage(message)

	requestBody, err := schema.GetIdentifyUser10RequestFromMessage(message)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.BAD_REQUEST)
		ms.bus.SendMessage(responseMessage)
		return
	}

	tokenPayload, err := ms.tokenCreator.Verify(requestBody.Token)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.NOT_AUTHORIZED)
		ms.bus.SendMessage(responseMessage)
		return
	}

	dbUser, err := ms.store.GetUserByUsername(context.Background(), tokenPayload.Username)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.NOT_AUTHORIZED)
		ms.bus.SendMessage(responseMessage)
		return
	}

	err = schema.FillBodyIdentifyUser10ResponseForMessage(dbUser.ID, &responseMessage)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.INTERNAL_ERROR)
		ms.bus.SendMessage(responseMessage)
		return
	}

	metaBody := schema.Meta{
		WsId:   message.Meta.WsId,
		UserId: dbUser.ID,
	}

	schema.FillMetaForMessage(&responseMessage, metaBody)

	schema.FillStatusForMessage(&responseMessage, schema.OK)
	ms.bus.SendMessage(responseMessage)
}

func (ms *Microservice) getDevelopers10(message schema.Message) {
	responseMessage := schema.GenerateResponseMessage(message)

	requestBody, err := schema.GetGetDevelopers10RequestFromMessage(message)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.BAD_REQUEST)
		ms.bus.SendMessage(responseMessage)
		return
	}

	dbRole, err := ms.store.GetRoleByName(context.Background(), ROLE_DEVELOPER)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.INTERNAL_ERROR)
		ms.bus.SendMessage(responseMessage)
		return
	}

	argsForFindDevelopers := database.GetUsersByRoleIdAndIdFromParams{
		RoleID: dbRole.ID,
		ID:     requestBody.FromUserId,
	}

	dbDevelopers, err := ms.store.GetUsersByRoleIdAndIdFrom(context.Background(), argsForFindDevelopers)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.INTERNAL_ERROR)
		ms.bus.SendMessage(responseMessage)
		return
	}

	developersBody := make([]schema.Developer, len(dbDevelopers))
	for i, value := range dbDevelopers {
		developersBody[i] = schema.GenerateDeveloperStruct(value.ID)
	}

	err = schema.FillBodyGetDevelopers10ResponseForMessage(developersBody, &responseMessage)
	if err != nil {
		schema.FillStatusForMessage(&responseMessage, schema.INTERNAL_ERROR)
		ms.bus.SendMessage(responseMessage)
		return
	}

	schema.FillStatusForMessage(&responseMessage, schema.OK)
	ms.bus.SendMessage(responseMessage)
}
