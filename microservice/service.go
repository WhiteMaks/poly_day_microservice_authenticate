package microservice

import (
	"github.com/gin-gonic/gin"
	"poly_day_microservice_authenticate/autogenerate/database"
)

const (
	PathCreateUser = "/user/create"
)

type Microservice struct {
	store  database.Store
	router *gin.Engine
}

func NewMicroservice(store database.Store) *Microservice {
	microservice := &Microservice{}

	router := gin.New()
	router.POST(PathCreateUser, microservice.createUser)

	microservice.store = store
	microservice.router = router

	return microservice
}

func (microservice *Microservice) Start(address string) error {
	return microservice.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
