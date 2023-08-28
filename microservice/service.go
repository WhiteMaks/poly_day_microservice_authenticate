package microservice

import (
	"github.com/gin-gonic/gin"
	"poly_day_microservice_authenticate/autogenerate/database"
)

type Service struct {
	store  *database.Store
	router *gin.Engine
}

func NewService(store *database.Store) *Service {
	service := &Service{}

	router := gin.Default()
	router.POST("/create_user", service.createUser)

	service.store = store
	service.router = router

	return service
}

func (service *Service) Start(address string) error {
	return service.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
