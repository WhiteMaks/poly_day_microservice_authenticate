package main

import (
	"database/sql"
	"github.com/WhiteMaks/poly_day_bus_schema/wbus"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/microservice"

	_ "github.com/lib/pq"
)

var (
	psqlConStr   = os.Getenv("POSTGRES_CONNECTING_STRING")
	rabbitConStr = os.Getenv("RABBIT_CONNECTING_STRING")
	appPort      = os.Getenv("APP_PORT")
	msName       = os.Getenv("MICROSERVICE_NAME")
)

func main() {
	bus := wbus.NewRabbit(msName)
	err := bus.Start(rabbitConStr)
	if err != nil {
		log.Fatalf("error start wbus: %v", err)
	}

	connection, err := sql.Open("postgres", psqlConStr)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	store := database.NewStore(connection)

	gin.SetMode(gin.ReleaseMode)

	ms := microservice.NewMs(store)
	err = ms.Start("0.0.0.0:" + appPort)
	if err != nil {
		log.Fatalf("error start microservice: %v", err)
	}
}
