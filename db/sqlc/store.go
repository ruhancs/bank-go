package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func Newstore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//parametro nil seria isolacao customizada para transacoes
	tx,err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//criar a transacao com db
	q := New(tx)
	err = fn(q) // fn executa a operacao do db
	if err != nil {
		//retorna a operacao
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb error: %v", err, rbErr)//erro no rolback
		}
		return err
	}
	return tx.Commit() // finaliza a transacao salvando no db
}

//parametros para transferir dinheiro entre as contas
type TranferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

//resuto da tranferencia entre contas
type TransferTxResult struct {
	Tranfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

//transfere dinheiro de uma conta para outra
// cria tranferencia na tabela transfers, cria entry na tabela entries, atualiza as contas de recebimento e de saida
//tudo em uma operacao com db
func (store *Store) TranferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Tranfer,err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount, // valor negativo pois o dinheiro esta saindo da conta
		})
		if err != nil {
			return err
		}
		
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount, // valor positivo pois o dinheiro esta entrando
		})
		if err != nil {
			return err
		}

		// get account => update balance
		
		result.FromAccount,err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount,err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result,err
}