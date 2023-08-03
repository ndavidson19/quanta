package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreateTx(ctx context.Context, arg CreateTxParams) (CreateTxResult, error)
}

// Store provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// This transaction performs the following operations atomically:
// 1. Create a new user
// 2. Create a new account for the user
// 3. Create a new entry in the accounts table representing the initial deposit
// 5. Create a new entry in the deposits table representing the initial deposit
// 6. Create a new entry in the audit_log table representing the initial deposit

type CreateTxParams struct {
	AccountID int64  `json:"account_id"`
	Owner     string `json:"owner"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}

type CreateTxResult struct {
	Account  Account  `json:"account"`
	Deposit  Deposit  `json:"deposit"`
	AuditLog AuditLog `json:"audit_log"`
	Balance  int64    `json:"balance"`
}

var txKey = struct{}{}

func (store *SQLStore) CreateTx(ctx context.Context, arg CreateTxParams) (CreateTxResult, error) {
	var result CreateTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(">> tx name: ", txName, "create account: ", arg.Owner, "amount: ", arg.Amount)

		result.Account, err = q.CreateAccount(ctx, CreateAccountParams{
			Owner:    arg.Owner,
			Balance:  arg.Amount,
			Currency: arg.Currency,
		})
		if err != nil {
			return err
		}

		fmt.Println(">> tx name: ", txName, "create deposit: ", arg.Owner, "amount: ", arg.Amount)

		result.Deposit, err = q.CreateDeposit(ctx, CreateDepositParams{
			AccountID: result.Account.ID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
