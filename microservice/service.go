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

func NewMs(store database.Store) *Microservice {
	ms := &Microservice{}

	router := gin.New()
	router.POST(PathCreateUser, ms.createUser)

	ms.store = store
	ms.router = router

	return ms
}

func (ms *Microservice) Start(address string) error {
	return ms.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
