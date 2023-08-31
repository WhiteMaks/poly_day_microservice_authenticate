package microservice

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/utils"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,max=25"`
	Password string `json:"password" binding:"required,max=72"`
}

type CreateUserResponse struct {
	Id int64 `json:"id"`
}

func (microservice *Microservice) createUser(ctx *gin.Context) {
	var requestBody CreateUserRequest

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExist, err := userByUsernameExistInDb(microservice, ctx, requestBody.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if userExist {
		errorMessage := errors.New("user with the same name already exists")
		ctx.JSON(http.StatusConflict, errorResponse(errorMessage))
		return
	}

	dbRole, err := microservice.store.GetRoleByName(ctx, "ROLE_DEVELOPER")
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

	dbUser, err := microservice.store.CreateUser(ctx, sqlUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := CreateUserResponse{
		Id: dbUser.ID,
	}

	ctx.JSON(http.StatusCreated, response)
}

func userByUsernameExistInDb(service *Microservice, ctx *gin.Context, username string) (bool, error) {
	_, err := service.store.GetUserByUsername(ctx, username)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	}

	return true, err
}
