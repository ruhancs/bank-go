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
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
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