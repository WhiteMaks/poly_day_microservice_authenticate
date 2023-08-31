package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	firstRandomPassword := RandomString(25)
	secondRandomPassword := RandomString(25)

	firstHashedPassword, err := HashPassword(firstRandomPassword)
	assert.NoError(t, err, "No error must be")
	assert.NotEmpty(t, firstHashedPassword, "Hashed password must not be empty")

	err = PasswordEquality(firstRandomPassword, firstHashedPassword)
	assert.NoError(t, err, "No error must be, because passwords are equals")

	err = PasswordEquality(secondRandomPassword, firstHashedPassword)
	assert.Error(t, err, "Error must be, because passwords are not equals")
}
