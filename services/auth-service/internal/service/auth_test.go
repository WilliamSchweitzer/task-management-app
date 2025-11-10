package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "SecurePassw0rd!"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestCheckPassword(t *testing.T) {
	password := "SecurePassw0rd!"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// Correct password
	isValid := CheckPassword(password, hashedPassword)
	assert.True(t, isValid)

	// Incorrect password
	isValid = CheckPassword("WrongPassword", hashedPassword)
	assert.False(t, isValid)
}
