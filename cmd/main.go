package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/microservice"

	_ "github.com/lib/pq"
)

var (
	psqlConStr = os.Getenv("POSTGRES_CONNECTING_STRING")
	appPort    = os.Getenv("APP_PORT")
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	connection, err := sql.Open("postgres", psqlConStr)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	store := database.NewStore(connection)
	ms := microservice.NewMs(store)

	err = ms.Start("0.0.0.0:" + appPort)
	if err != nil {
		log.Fatalf("error start microservice: %v", err)
	}
}
