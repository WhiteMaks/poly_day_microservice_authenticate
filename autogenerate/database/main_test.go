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

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var testQueries *Queries
var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randomString(characters int) string {
	b := make([]rune, characters)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	connection, err := sql.Open(DbDriver, DbSource)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	testQueries = New(connection)

	os.Exit(m.Run())
}
