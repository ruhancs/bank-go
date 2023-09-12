package db

import (
	"context"
	"testing"

	"github.com/ruhancs/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry,err := testQueries.CreateEntry(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,entry)
	require.Equal(t,entry.Amount,arg.Amount)
	require.Equal(t,entry.AccountID,arg.AccountID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t,account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t,account)

	foundedEntry,err := testQueries.GetEntry(context.Background(),entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, foundedEntry)
	require.Equal(t,entry.ID,foundedEntry.ID)
	require.Equal(t,entry.AccountID,foundedEntry.AccountID)
	require.Equal(t,entry.Amount,foundedEntry.Amount)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t,account)
	}

	arg := ListEntryParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}
	entries,err := testQueries.ListEntry(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,entries,5)

	for _,entry := range entries {
		require.NotEmpty(t,entry)
	}
}