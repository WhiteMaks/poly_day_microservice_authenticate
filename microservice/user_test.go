package microservice

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/utils"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DbDriver = "postgres"
	DbSource = "postgres://user:password@localhost:15432/authenticate?sslmode=disable"
)

func TestMicroservice_CreateUser(t *testing.T) {
	connection, err := sql.Open(DbDriver, DbSource)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	store := database.NewStore(connection)
	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateRandomMSCreateUserRequest()
	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusCreated, recorder.Code, "response [ status code ] must be equal")

	dbUser, err := store.GetUserByUsername(context.Background(), createUserBody.Username)
	assert.NoError(t, err, "No error must be")

	checkResponseCreateUser(t, createUserBody, recorder.Body, dbUser)
}

func checkResponseCreateUser(t *testing.T, requestBody createUserRequest, responseBody *bytes.Buffer, dbUser database.RootUser) {
	data, err := io.ReadAll(responseBody)
	assert.NoError(t, err, "No error must be")

	var userResponse createUserResponse
	err = json.Unmarshal(data, &userResponse)
	assert.NoError(t, err, "No error must be")

	assert.Equal(t, userResponse.Id, dbUser.ID, "user [ id ] must be equal")
	assert.Equal(t, requestBody.Username, dbUser.Username, "user [ username ] must be equal")
	assert.Equal(t, requestBody.Password, dbUser.Password, "user [ password ] must be equal")
}

func generateRandomMSCreateUserRequest() createUserRequest {
	return createUserRequest{
		Username: utils.RandomString(25),
		Password: utils.RandomString(256),
	}
}
