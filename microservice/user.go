package microservice

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"poly_day_microservice_authenticate/autogenerate/database"
)

type createUserRequest struct {
	Username string `json:"username" biding:"required"`
	Password string `json:"password" biding:"required"`
}

type createUserResponse struct {
	Id int64 `json:"id"`
}

func (microservice *Microservice) createUser(ctx *gin.Context) {
	var requestBody createUserRequest

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExist := userByUsernameExistInDb(microservice, ctx, requestBody.Username)
	if userExist {
		ctx.JSON(http.StatusConflict, errorResponse(errors.New("user with the same name already exists")))
		return
	}

	dbRole, err := microservice.store.GetRoleByName(ctx, "ROLE_DEVELOPER")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sqlUser := database.CreateUserParams{
		Username: requestBody.Username,
		Password: requestBody.Password,
		RoleID:   dbRole.ID,
	}

	dbUser, err := microservice.store.CreateUser(ctx, sqlUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse{
		Id: dbUser.ID,
	}

	ctx.JSON(http.StatusCreated, response)
}

func userByUsernameExistInDb(service *Microservice, ctx *gin.Context, username string) bool {
	_, err := service.store.GetUserByUsername(ctx, username)
	if err != nil {
		return false
	}

	return true
}
