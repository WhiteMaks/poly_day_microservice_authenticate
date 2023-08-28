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
	id int64
}

func (service *Service) createUser(ctx *gin.Context) {
	var requestBody createUserRequest

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExist := userByUsernameExistInDb(service, ctx, requestBody.Username)
	if userExist {
		ctx.JSON(http.StatusConflict, errorResponse(errors.New("user with the same name already exists")))
		return
	}

	dbRole, err := service.store.GetRoleByName(ctx, "ROLE_DEVELOPER")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sqlUser := database.CreateUserParams{
		Username: requestBody.Username,
		Password: requestBody.Password,
		RoleID:   dbRole.ID,
	}

	dbUser, err := service.store.CreateUser(ctx, sqlUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse{
		id: dbUser.ID,
	}

	ctx.JSON(http.StatusCreated, response)
}

func userByUsernameExistInDb(service *Service, ctx *gin.Context, username string) bool {
	_, err := service.store.GetUserByUsername(ctx, username)
	if err != nil {
		return false
	}

	return true
}
