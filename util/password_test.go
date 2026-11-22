package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	wrongPassword := RandomString(10)
	require.NoError(t, err)

	err = CheckPassword(wrongPassword, hashed)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	password = RandomString(10)
	hashed1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed1)

	require.NotEqual(t, hashed, hashed1)
}