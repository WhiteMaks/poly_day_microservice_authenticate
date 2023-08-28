package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueries_CreateUser(t *testing.T) {
	randomUser := generateRandomUser()

	actualUser, err := testQueries.CreateUser(context.Background(), randomUser)
	assert.NoError(t, err, "No error must be")

	expectedUser := RootUser{
		Username: randomUser.Username,
		Password: randomUser.Password,
		RoleID:   randomUser.RoleID,
	}

	checkUsers(t, expectedUser, actualUser)
}

func TestQueries_CreateUser_WithNotUniqueName(t *testing.T) {
	existUser := generateRandomUser()

	_, err := testQueries.CreateUser(context.Background(), existUser)
	assert.NoError(t, err, "No error must be")

	_, err = testQueries.CreateUser(context.Background(), existUser)
	assert.Error(t, err, "Error must be")
	assert.Equal(t, "pq: duplicate key value violates unique constraint \"users_username_key\"", err.Error())
}

func TestQueries_GetUserByUsername(t *testing.T) {
	randomUser := generateRandomUser()

	expectedUser, err := testQueries.CreateUser(context.Background(), randomUser)
	assert.NoError(t, err, "No error must be")

	actualUser, err := testQueries.GetUserByUsername(context.Background(), randomUser.Username)
	assert.NoError(t, err, "No error must be")

	checkUsers(t, expectedUser, actualUser)
}

func TestQueries_GetUserByUsername_WithNotExistValue(t *testing.T) {
	randomUser := generateRandomUser()

	notExistUser, err := testQueries.GetUserByUsername(context.Background(), randomUser.Username)
	assert.Error(t, err, "Error must be")
	assert.Equal(t, "sql: no rows in result set", err.Error())

	assert.Empty(t, notExistUser, "User must be empty")
}

func checkUsers(t *testing.T, expectedUser RootUser, actualUser RootUser) {
	assert.NotZero(t, actualUser.ID, "User [ id ] must be filled")
	assert.NotZero(t, actualUser.CreatedAt, "User [ created_at ] must be filled")
	assert.Equal(t, expectedUser.Username, actualUser.Username, "User [ username ] must be equal")
	assert.Equal(t, expectedUser.Password, actualUser.Password, "User [ password ] must be equal")
	assert.Equal(t, expectedUser.RoleID, actualUser.RoleID, "User [ role_id ] must be equal")
}

func generateRandomUser() CreateUserParams {
	return CreateUserParams{
		Username: randomString(25),
		Password: randomString(256),
		RoleID:   expectedRoles[1].ID,
	}
}
