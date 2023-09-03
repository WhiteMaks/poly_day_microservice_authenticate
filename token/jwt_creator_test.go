package token

import (
	"github.com/stretchr/testify/assert"
	"poly_day_microservice_authenticate/utils"
	"testing"
	"time"
)

func TestJWTCreator_ValidToken(t *testing.T) {
	creator := NewJWTCreator(utils.RandomString(10))

	expectedUsername := utils.RandomString(25)
	duration := time.Minute

	token, err := creator.Create(expectedUsername, duration)
	assert.NoError(t, err, "No error must be")
	assert.NotEmpty(t, token, "Token must not be empty")

	payload, err := creator.Verify(token)
	assert.NoError(t, err, "No error must be")
	assert.NotNil(t, payload, "Payload must not be nil")

	assert.NotEmpty(t, payload.CreateAt, "CreateAt must not be empty")
	assert.NotEmpty(t, payload.ExpiredAt, "CreateAt must not be empty")
	assert.Equal(t, expectedUsername, payload.Username, "Username must be equals")
}

func TestJWTCreator_ExpiredToken(t *testing.T) {
	creator := NewJWTCreator(utils.RandomString(10))

	expectedUsername := utils.RandomString(25)
	duration := -time.Minute

	token, err := creator.Create(expectedUsername, duration)
	assert.NoError(t, err, "No error must be")
	assert.NotEmpty(t, token, "Token must not be empty")

	payload, err := creator.Verify(token)
	assert.Error(t, err, "No error must be")
	assert.EqualError(t, err, ErrExpiredToken.Error(), "Error message must be equals")
	assert.Nil(t, payload, "Payload must be nil")
}
