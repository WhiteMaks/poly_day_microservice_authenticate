package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/microservice"

	_ "github.com/lib/pq"
)

const (
	DbDriver = "postgres"
	DbSource = "postgres://user:password@localhost:15432/authenticate?sslmode=disable"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	connection, err := sql.Open(DbDriver, DbSource)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	store := database.NewStore(connection)
	ms := microservice.NewMicroservice(store)

	err = ms.Start("localhost:8080")
	if err != nil {
		log.Fatalf("error start microservice: %v", err)
	}
}
