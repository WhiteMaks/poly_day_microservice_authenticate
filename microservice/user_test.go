package microservice

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"poly_day_microservice_authenticate/autogenerate/database"
	"poly_day_microservice_authenticate/autogenerate/database/mock"
	"poly_day_microservice_authenticate/utils"
	"testing"
	"time"
)

var developerRole = database.RootRole{
	ID:   1,
	Name: "ROLE_DEVELOPER",
}

func TestMicroservice_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbUser := generateRandomDBUser()
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserMockDB(store, dbUser)

	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserResponseFromDBRecord(dbUser)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusCreated, recorder.Code, "response [ status code ] must be equal")

	checkResponseCreateUser(t, createUserBody, recorder.Body, dbUser)
}

func TestMicroservice_CreateUser_InternalError_ErrorInGetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbUser := generateRandomDBUser()
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorGetUserByUsernameMockDB(store, dbUser)

	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserResponseFromDBRecord(dbUser)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_CreateUser_InternalError_ErrorInGetRoleByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbUser := generateRandomDBUser()
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorGetRoleByNameMockDB(store, dbUser)

	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserResponseFromDBRecord(dbUser)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_CreateUser_InternalError_ErrorInCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbUser := generateRandomDBUser()
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorCreateUserMockDB(store, dbUser)

	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserResponseFromDBRecord(dbUser)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_CreateUser_BudRequest_WrongJsonBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdatabase.NewMockStore(ctrl)
	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	requests := []struct {
		request []byte
	}{
		{request: []byte{}},
		{request: []byte("{\"username\":\"123\"}")},
		{request: []byte("{\"password\":\"123\"}")},
		{request: []byte("{\"passwor\":\"123\", \"usernam\":\"123\"}")},
		{request: []byte("{\"password\":\"123\", \"username\":\"12345678901234567890123456\"}")},
	}

	for _, requestStruct := range requests {
		request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestStruct.request))
		assert.NoError(t, err, "No error must be")

		ms.router.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "response [ status code ] must be equal")
	}
}

func TestMicroservice_CreateUser_Conflict_UserExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbUser := generateRandomDBUser()
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserConflictUserExistMockDB(store, dbUser)

	ms := NewMicroservice(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserResponseFromDBRecord(dbUser)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusConflict, recorder.Code, "response [ status code ] must be equal")
}

func testMicroserviceCreateUserMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, sql.ErrNoRows)

	store.EXPECT().
		GetRoleByName(
			gomock.Any(),
			gomock.Eq(developerRole.Name),
		).Times(1).
		Return(developerRole, nil)

	store.EXPECT().
		CreateUser(
			gomock.Any(),
			gomock.Eq(
				database.CreateUserParams{
					Username: dbUser.Username,
					Password: dbUser.Password,
					RoleID:   dbUser.RoleID,
				},
			),
		).Times(1).
		Return(dbUser, nil)
}

func testMicroserviceCreateUserInternalErrorErrorGetUserByUsernameMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, errors.New("unexpected error"))
}

func testMicroserviceCreateUserInternalErrorErrorGetRoleByNameMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, sql.ErrNoRows)

	store.EXPECT().
		GetRoleByName(
			gomock.Any(),
			gomock.Eq(developerRole.Name),
		).Times(1).
		Return(database.RootRole{}, errors.New("unexpected error"))
}

func testMicroserviceCreateUserInternalErrorErrorCreateUserMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, sql.ErrNoRows)

	store.EXPECT().
		GetRoleByName(
			gomock.Any(),
			gomock.Eq(developerRole.Name),
		).Times(1).
		Return(developerRole, nil)

	store.EXPECT().
		CreateUser(
			gomock.Any(),
			gomock.Eq(
				database.CreateUserParams{
					Username: dbUser.Username,
					Password: dbUser.Password,
					RoleID:   dbUser.RoleID,
				},
			),
		).Times(1).
		Return(database.RootUser{}, errors.New("unexpected error"))
}

func testMicroserviceCreateUserConflictUserExistMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, nil)
}

func checkResponseCreateUser(t *testing.T, requestBody CreateUserRequest, responseBody *bytes.Buffer, dbUser database.RootUser) {
	data, err := io.ReadAll(responseBody)
	assert.NoError(t, err, "No error must be")

	var userResponse CreateUserResponse
	err = json.Unmarshal(data, &userResponse)
	assert.NoError(t, err, "No error must be")

	assert.Equal(t, userResponse.Id, dbUser.ID, "user [ id ] must be equal")
	assert.Equal(t, requestBody.Username, dbUser.Username, "user [ username ] must be equal")
	assert.Equal(t, requestBody.Password, dbUser.Password, "user [ password ] must be equal")
}

func generateCreateUserResponseFromDBRecord(dbUser database.RootUser) CreateUserRequest {
	return CreateUserRequest{
		Username: dbUser.Username,
		Password: dbUser.Password,
	}
}

func generateRandomDBUser() database.RootUser {
	return database.RootUser{
		ID:        utils.RandomInt(0, 10000),
		Username:  utils.RandomString(25),
		Password:  utils.RandomString(256),
		RoleID:    developerRole.ID,
		CreatedAt: time.Now(),
	}
}
