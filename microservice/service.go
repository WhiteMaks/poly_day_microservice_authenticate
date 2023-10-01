package microservice

import (
	"github.com/WhiteMaks/poly_day_bus_schema/schema"
	"github.com/WhiteMaks/poly_day_bus_schema/wbus"
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
	bus          *wbus.Rabbit
}

func NewMs(store database.Store, bus *wbus.Rabbit) *Microservice {
	ms := &Microservice{}

	router := gin.New()

	router.GET(PathCheckUser, ms.checkUser)

	router.POST(PathCreateUser, ms.createUser)
	router.POST(PathLoginUser, ms.loginUser)

	ms.bus = bus
	ms.store = store
	ms.router = router
	ms.tokenCreator = token.NewJWTCreator("secret_key_for_generate_token") //todo получать секретный ключ из конфига

	return ms
}

func (ms *Microservice) Start(address string) error {
	err := ms.bus.SubscribeOnQueue(schema.PathMsAuthenticate, ms.rootMessageHandler)
	if err != nil {
		return err
	}

	err = ms.router.Run(address)
	if err != nil {
		return err
	}

	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (ms *Microservice) rootMessageHandler(message schema.Message) {
	methodName := generateFullMethodNameFromMessage(message)

	if methodName == schema.CommandIdentifyUser10 {
		go ms.identifyUser10(message)
		return
	}

}

func generateFullMethodNameFromMessage(message schema.Message) string { //todo вынести метод в bus schema
	return message.Method + "-" + message.Version
}
