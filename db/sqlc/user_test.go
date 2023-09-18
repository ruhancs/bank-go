package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ruhancs/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPasswod, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams {
		Username: util.RandomOwner(),
		HashedPassword: hashedPasswod,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	//testQueries declarado em main_test.go
	user,err := testQueries.CreateUser(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t,user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.True(t,user.PasswordChangedAt.IsZero())
	//require.NotZero(t,user.PasswordChangedAt)
	require.NotZero(t,user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	foundedUser,err := testQueries.GetUser(context.Background(), createdUser.Username)

	require.NoError(t,err)
	require.Equal(t,foundedUser.Username,createdUser.Username)
	require.Equal(t,foundedUser.HashedPassword,createdUser.HashedPassword)
	require.Equal(t,foundedUser.FullName,createdUser.FullName)
	require.WithinDuration(t, foundedUser.CreatedAt,createdUser.CreatedAt, time.Second)
}

func TestUpdateUserFullName(t *testing.T) {
	createdUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	updatedUser,err := testQueries.UpdateUser(context.Background(),UpdateUserParams{
		Username: createdUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
	})

	require.NoError(t,err)
	require.Equal(t,updatedUser.FullName, newFullName)
	require.Equal(t,updatedUser.Email,createdUser.Email)
	require.Equal(t,updatedUser.HashedPassword,createdUser.HashedPassword)
}

func TestUpdateUserEmail(t *testing.T) {
	createdUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser,err := testQueries.UpdateUser(context.Background(),UpdateUserParams{
		Username: createdUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid: true,
		},
	})

	require.NoError(t,err)
	require.Equal(t,updatedUser.Email, newEmail)
	require.Equal(t,updatedUser.FullName,createdUser.FullName)
	require.Equal(t,updatedUser.HashedPassword,createdUser.HashedPassword)
}