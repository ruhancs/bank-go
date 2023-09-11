package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := Newstore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("Beffore",account1.Balance, account2.Balance)
	fmt.Println(account2.Balance)

	// rodar n operacoes de transferencia e concorrente
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	existed := make(map[int]bool)
	//verificar os resultados
	for i := 0; i < n; i++ {
		
		go func() {
			//contexto para inserir o nome da transacao, txKey criado em store.go
			ctx := context.Background()
			result,err := store.TranferTx(ctx,TranferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t,err)

		result := <-results
		require.NotEmpty(t,result)
		
		transfer := result.Tranfer
		require.NotEmpty(t,transfer)
		require.Equal(t,transfer.FromAccountID,account1.ID)
		require.Equal(t,transfer.ToAccountID,account2.ID)
		require.Equal(t,transfer.Amount,amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)

		_,err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t,err)
		
		fromEntry := result.FromEntry
		require.NotEmpty(t,fromEntry)
		require.Equal(t,account1.ID,fromEntry.AccountID)
		require.Equal(t,-amount,fromEntry.Amount)
		require.NotZero(t,fromEntry.ID)
		require.NotZero(t,fromEntry.CreatedAt)
		
		_,err = store.GetEntry(context.Background(),fromEntry.ID)
		require.NoError(t,err)
		
		toEntry := result.ToEntry
		require.NotEmpty(t,toEntry)
		require.Equal(t,account2.ID,toEntry.AccountID)
		require.Equal(t,amount,toEntry.Amount)
		require.NotZero(t,toEntry.ID)
		require.NotZero(t,toEntry.CreatedAt)
		
		_,err = store.GetEntry(context.Background(),toEntry.ID)
		require.NoError(t,err)

		fromAccount := result.FromAccount
		require.NotEmpty(t,fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		
		toAccount := result.ToAccount
		require.NotEmpty(t,toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t,diff1,diff2)
		require.True(t,diff1 > 0)
		require.True(t,diff1 % amount == 0)

		k:= int(diff1/amount)
		require.True(t,k >=1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccont1,err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t,err)
	
	updatedAccont2,err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t,err)

	fmt.Println("Beffore",updatedAccont1.Balance, updatedAccont2.Balance)
	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccont1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccont2.Balance)

}