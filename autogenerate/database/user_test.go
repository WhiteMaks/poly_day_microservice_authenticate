package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueries_CreateUser(t *testing.T) {
	expectedUser := generateRandomUser()

	actualUser, err := testQueries.CreateUser(context.Background(), expectedUser)
	assert.NoError(t, err, "No error must be")

	assert.NotZero(t, actualUser.ID, "User [ id ] must be filled")
	assert.NotZero(t, actualUser.CreatedAt, "User [ created_at ] must be filled")
	assert.Equal(t, expectedUser.Username, actualUser.Username, "User [ username ] must be equal")
	assert.Equal(t, expectedUser.Password, actualUser.Password, "User [ password ] must be equal")
	assert.Equal(t, expectedUser.RoleID, actualUser.RoleID, "User [ role_id ] must be equal")
}

func TestQueries_CreateUser_WithNotUniqueName(t *testing.T) {
	existUser := generateRandomUser()

	_, err := testQueries.CreateUser(context.Background(), existUser)
	assert.NoError(t, err, "No error must be")

	_, err = testQueries.CreateUser(context.Background(), existUser)
	assert.Error(t, err, "Error must be")
	assert.Equal(t, "pq: duplicate key value violates unique constraint \"users_username_key\"", err.Error())
}

func generateRandomUser() CreateUserParams {
	return CreateUserParams{
		Username: randomString(25),
		Password: randomString(256),
		RoleID:   expectedRoles[1].ID,
	}
}
