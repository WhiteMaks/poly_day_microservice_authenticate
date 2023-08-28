package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedRoles = []RootRole{
	{
		ID:   1,
		Name: "ROLE_DEVELOPER",
	},
	{
		ID:   2,
		Name: "ROLE_GAMER",
	},
}

func TestQueries_GetAllRoles(t *testing.T) {
	actualRoles, err := testQueries.GetAllRoles(context.Background())
	assert.NoError(t, err, "No error must be")

	assert.Equal(t, len(expectedRoles), len(actualRoles), "Rows size")
	assert.Equal(t, expectedRoles, actualRoles, "Roles must be equals")
}

func TestQueries_GetRoleById(t *testing.T) {
	for _, expectedRole := range expectedRoles {
		actualRole, err := testQueries.GetRoleById(context.Background(), expectedRole.ID)
		assert.NoError(t, err, "No error must be")

		assert.NotEmpty(t, actualRole, "Response must be not empty")
		assert.Equal(t, expectedRole.ID, actualRole.ID, "Role [ id ] must be equal")
		assert.Equal(t, expectedRole.Name, actualRole.Name, "Role [ name ] must be equal")
	}
}

func TestQueries_GetRoleByName(t *testing.T) {
	for _, expectedRole := range expectedRoles {
		actualRole, err := testQueries.GetRoleByName(context.Background(), expectedRole.Name)
		assert.NoError(t, err, "No error must be")

		assert.NotEmpty(t, actualRole, "Response must be not empty")
		assert.Equal(t, expectedRole.ID, actualRole.ID, "Role [ id ] must be equal")
		assert.Equal(t, expectedRole.Name, actualRole.Name, "Role [ name ] must be equal")
	}
}
