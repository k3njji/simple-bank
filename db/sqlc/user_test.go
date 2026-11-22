package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashed_password, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	
	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		HashedPassword:  hashed_password,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, arg.Username, user.Username)
	require.NotZero(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, arg.FullName, user.FullName)
	require.NotZero(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	// require.NotZero(t, user.ID)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	account1 :=  createRandomUser(t)

	account2, err := testQueries.GetUser(context.Background(), account1.Username)

	require.NoError(t, err)

	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Email, account2.Email)
}

// func TestUpdateAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	arg := UpdateAccountParams{
// 		ID:      account1.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	account2, err := testQueries.UpdateAccount(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)
// 	require.Equal(t, arg.ID, account2.ID)
// 	require.Equal(t, arg.Balance, account2.Balance)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.NotEqual(t, account1.Balance, account2.Balance)
// }

// func TestDeleteAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.Error(t, err)
// 	require.Empty(t, account2)
// }

// func TestListAccounts(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}

// 	arg := ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccounts(context.Background(), arg)

// 	require.NoError(t, err)
// 	for _, accounts := range accounts {
// 		require.NotEmpty(t, accounts)
// 	}
// }
