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

func TestMicroservice_CheckUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdatabase.NewMockStore(ctrl)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	expectedUsername := utils.RandomString(25)
	duration := time.Minute

	generatedToken, err := ms.tokenCreator.Create(expectedUsername, duration)
	assert.NoError(t, err, "No error must be")
	assert.NotEmpty(t, generatedToken, "Token must not be empty")

	request, err := http.NewRequest(http.MethodGet, PathCheckUser, nil)
	assert.NoError(t, err, "No error must be")

	request.Header.Add("Authorization", generatedToken)

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "response [ status code ] must be equal")

	checkResponseCheckUser(t, recorder.Body, true)
}

func TestMicroservice_CheckUser_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdatabase.NewMockStore(ctrl)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, PathCheckUser, nil)
	assert.NoError(t, err, "No error must be")

	request.Header.Add("Authorization", "invalid_token")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "response [ status code ] must be equal")

	checkResponseCheckUser(t, recorder.Body, false)
}

func TestMicroservice_CheckUser_EmptyAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdatabase.NewMockStore(ctrl)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, PathCheckUser, nil)
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "response [ status code ] must be equal")

	checkResponseCheckUser(t, recorder.Body, false)
}

func TestMicroservice_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	loginUserBody := generateLoginUserRequestFromDBRecord(dbUser, randomPassword)
	store := mockdatabase.NewMockStore(ctrl)

	testMicroserviceLoginUserMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requestBody, err := json.Marshal(loginUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathLoginUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "response [ status code ] must be equal")

	checkResponseLoginUser(t, recorder.Body)
}

func TestMicroservice_LoginUser_Unauthorized_WithNotValidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	loginUserBody := generateLoginUserRequestFromDBRecord(dbUser, "not_valid_password")
	store := mockdatabase.NewMockStore(ctrl)

	testMicroserviceLoginUserMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requestBody, err := json.Marshal(loginUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathLoginUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_LoginUser_Unauthorized_WithNotValidUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	loginUserBody := generateLoginUserRequestFromDBRecord(dbUser, randomPassword)
	store := mockdatabase.NewMockStore(ctrl)

	testMicroserviceLoginUserUnauthorizedErrorUserNotExistMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requestBody, err := json.Marshal(loginUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathLoginUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_LoginUser_InternalError_ErrorInGetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	loginUserBody := generateLoginUserRequestFromDBRecord(dbUser, randomPassword)
	store := mockdatabase.NewMockStore(ctrl)

	testMicroserviceLoginUserInternalErrorErrorGetUserByUsernameMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requestBody, err := json.Marshal(loginUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathLoginUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "response [ status code ] must be equal")
}

func TestMicroservice_LoginUser_BudRequest_WrongJsonBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdatabase.NewMockStore(ctrl)
	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requests := []struct {
		request []byte
	}{
		{request: []byte{}},
		{request: []byte("{\"username\":\"123\"}")},
		{request: []byte("{\"password\":\"123\"}")},
		{request: []byte("{\"passwor\":\"123\", \"usernam\":\"123\"}")},
		{request: []byte("{\"passwor\":\"123\"; \"usernam\":\"123\"}")},
		{request: []byte("{\"password\":\"123\", \"username\":\"\"}")},
		{request: []byte("{\"password\":\"\", \"username\":\"123\"}")},
	}

	for _, requestStruct := range requests {
		request, err := http.NewRequest(http.MethodPost, PathLoginUser, bytes.NewBuffer(requestStruct.request))
		assert.NoError(t, err, "No error must be")

		ms.router.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "response [ status code ] must be equal")
	}
}

func TestMicroservice_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	createUserBody := generateCreateUserRequestFromDBRecord(dbUser, randomPassword)
	store := mockdatabase.NewMockStore(ctrl)

	testMicroserviceCreateUserMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

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

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorGetUserByUsernameMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserRequestFromDBRecord(dbUser, randomPassword)

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

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorGetRoleByNameMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserRequestFromDBRecord(dbUser, randomPassword)

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

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserInternalErrorErrorCreateUserMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserRequestFromDBRecord(dbUser, randomPassword)

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
	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	requests := []struct {
		request []byte
	}{
		{request: []byte{}},
		{request: []byte("{\"username\":\"123\"}")},
		{request: []byte("{\"password\":\"123\"}")},
		{request: []byte("{\"passwor\":\"123\", \"usernam\":\"123\"}")},
		{request: []byte("{\"passwor\":\"123\"; \"usernam\":\"123\"}")},
		{request: []byte("{\"password\":\"123\", \"username\":\"12345678901234567890123456\"}")},
		{request: []byte("{\"password\":\"123\", \"username\":\"\"}")},
		{request: []byte("{\"password\":\"\", \"username\":\"123\"}")},
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

	randomPassword := utils.RandomString(72)
	dbUser := generateRandomDBUser(randomPassword)
	store := mockdatabase.NewMockStore(ctrl)
	testMicroserviceCreateUserConflictUserExistMockDB(store, dbUser)

	ms := NewMs(store)
	recorder := httptest.NewRecorder()

	createUserBody := generateCreateUserRequestFromDBRecord(dbUser, randomPassword)

	requestBody, err := json.Marshal(createUserBody)
	assert.NoError(t, err, "No error must be")

	request, err := http.NewRequest(http.MethodPost, PathCreateUser, bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "No error must be")

	ms.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusConflict, recorder.Code, "response [ status code ] must be equal")
}

func testMicroserviceLoginUserMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(dbUser, nil)
}

func testMicroserviceLoginUserUnauthorizedErrorUserNotExistMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, sql.ErrNoRows)
}
func testMicroserviceLoginUserInternalErrorErrorGetUserByUsernameMockDB(store *mockdatabase.MockStore, dbUser database.RootUser) {
	store.EXPECT().
		GetUserByUsername(
			gomock.Any(),
			gomock.Eq(dbUser.Username),
		).Times(1).
		Return(database.RootUser{}, errors.New("unexpected error"))
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
			gomock.Any(),
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
			gomock.Any(),
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
	assert.NoError(t, utils.PasswordEquality(requestBody.Password, dbUser.Password), "user [ password ] must be equal")
}

func checkResponseLoginUser(t *testing.T, responseBody *bytes.Buffer) {
	data, err := io.ReadAll(responseBody)
	assert.NoError(t, err, "No error must be")

	var loginResponse LoginUserResponse
	err = json.Unmarshal(data, &loginResponse)
	assert.NoError(t, err, "No error must be")

	assert.NotEmpty(t, loginResponse.Token, "Token must not be empty")
}

func checkResponseCheckUser(t *testing.T, responseBody *bytes.Buffer, authorized bool) {
	data, err := io.ReadAll(responseBody)
	assert.NoError(t, err, "No error must be")

	var checkResponse CheckUserResponse
	err = json.Unmarshal(data, &checkResponse)
	assert.NoError(t, err, "No error must be")

	assert.Equal(t, authorized, checkResponse.Authorized, "Authorized must be equals")
}

func generateCreateUserRequestFromDBRecord(dbUser database.RootUser, password string) CreateUserRequest {
	return CreateUserRequest{
		Username: dbUser.Username,
		Password: password,
	}
}

func generateLoginUserRequestFromDBRecord(dbUser database.RootUser, password string) LoginUserRequest {
	return LoginUserRequest{
		Username: dbUser.Username,
		Password: password,
	}
}

func generateRandomDBUser(password string) database.RootUser {
	hashedPassword, _ := utils.HashPassword(password)

	return database.RootUser{
		ID:        utils.RandomInt(0, 10000),
		Username:  utils.RandomString(25),
		Password:  hashedPassword,
		RoleID:    developerRole.ID,
		CreatedAt: time.Now(),
	}
}
