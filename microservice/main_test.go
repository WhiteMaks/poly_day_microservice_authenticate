package microservice

import (
	"github.com/WhiteMaks/poly_day_bus_schema/schema"
	"github.com/WhiteMaks/poly_day_bus_schema/wbus"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"poly_day_microservice_authenticate/autogenerate/database/mock"
	"testing"
)

const (
	rabbitConStr = "amqp://user:password@localhost:5672/"
)

var bus *wbus.Rabbit

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)

	bus = wbus.NewRabbit(schema.PathMsAuthenticate)
	err := bus.Start(rabbitConStr)
	if err != nil {
		log.Fatalf("error start bus [ %v ]: %v", bus.Name, err)
	}

	os.Exit(m.Run())
}

func newMs(store *mockdatabase.MockStore) *Microservice {
	return NewMs(store, bus)
}
