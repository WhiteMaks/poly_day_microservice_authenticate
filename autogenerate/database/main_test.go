package database

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	DbDriver = "postgres"
	DbSource = "postgres://user:password@localhost:15432/authenticate?sslmode=disable"
)

var testQueries *Queries
var _ *rand.Rand

func init() {
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestMain(m *testing.M) {
	connection, err := sql.Open(DbDriver, DbSource)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	testQueries = New(connection)

	os.Exit(m.Run())
}
