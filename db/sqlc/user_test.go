package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/timly278/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	hashedPass, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	User, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, User)
	require.Equal(t, arg.Username, User.Username)
	require.Equal(t, arg.HashedPassword, User.HashedPassword)
	require.Equal(t, arg.FullName, User.FullName)
	require.Equal(t, arg.Email, User.Email)

	require.True(t, User.PasswordChangedAt.IsZero())
	require.NotZero(t, User.CreatedAt)

	return User
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	User1 := createRandomUser(t)
	User2, err := testQueries.GetUser(context.Background(), User1.Username)

	require.NoError(t, err)
	require.Equal(t, User1.Username, User2.Username)
	require.Equal(t, User1.HashedPassword, User2.HashedPassword)
	require.Equal(t, User1.FullName, User2.FullName)
	require.Equal(t, User1.Email, User2.Email)
	require.WithinDuration(t, User1.CreatedAt, User2.CreatedAt, time.Second)
}
