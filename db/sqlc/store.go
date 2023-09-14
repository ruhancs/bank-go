package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TranferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error)
	execTx(ctx context.Context, fn func(*Queries) error) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func Newstore(db *sql.DB) Store{
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
func (store *SQLStore) TranferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error) {
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

		//evitar deadlock de transacoes cruzadas
		//adicionar ordem nas atualizacoes das contas
		if arg.FromAccountID < arg.ToAccountID {
			// get account => update balance
			result.FromAccount,result.ToAccount,err = addMOney(ctx,q,arg.FromAccountID,-arg.Amount,arg.ToAccountID,arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount,result.FromAccount,err = addMOney(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result,err
}

func addMOney(
	ctx context.Context, 
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	
	account2,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	//retorna automatico account1,account2,err
	return 
}