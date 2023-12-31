package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ruhancs/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	//testQueries declarado em main_test.go
	account,err := testQueries.CreateAccount(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t,account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	foundedAccount,err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t,err)
	require.Equal(t,foundedAccount.Owner,createdAccount.Owner)
	require.Equal(t,foundedAccount.Balance,createdAccount.Balance)
	require.Equal(t,foundedAccount.Currency,createdAccount.Currency)
	require.WithinDuration(t, foundedAccount.CreatedAt,createdAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountParams {
		ID: createdAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount,err := testQueries.UpdateAccount(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,updatedAccount)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(),createdAccount.ID)

	require.NoError(t,err)

	account,err := testQueries.GetAccount(context.Background(),createdAccount.ID)
	require.Error(t,err)
	require.EqualError(t,err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams {
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accounts,err := testQueries.ListAccounts(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,accounts)

	for _, account := range accounts {
		require.NotEmpty(t,account)
	}
}

