package microservice

import (
	"github.com/gin-gonic/gin"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/token"
)

const (
	PathCreateUser = "/user/create"
	PathLoginUser  = "/user/login"
	PathCheckUser  = "/user/check"
)

type Microservice struct {
	store        database.Store
	tokenCreator token.Creator
	router       *gin.Engine
}

func NewMs(store database.Store) *Microservice {
	ms := &Microservice{}

	router := gin.New()

	router.GET(PathCheckUser, ms.checkUser)

	router.Use(CORSMiddleware())
	router.POST(PathCreateUser, ms.createUser)
	router.POST(PathLoginUser, ms.loginUser)

	ms.store = store
	ms.router = router
	ms.tokenCreator = token.NewJWTCreator("secret_key_for_generate_token")

	return ms
}

func (ms *Microservice) Start(address string) error {
	return ms.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
