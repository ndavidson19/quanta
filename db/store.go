package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type CreateTxResult struct {
	Account  Account  `json:"account"`
	Deposit  Deposit  `json:"deposit"`
	AuditLog AuditLog `json:"audit_log"`
	Balance  int64    `json:"balance"`
}

var txKey = struct{}{}

func (store *Store) CreateTx(ctx context.Context, arg CreateTxParams) (CreateTxResult, error) {
	var result CreateTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(">> tx name: ", txName, "create account: ", arg.AccountID, "amount: ", arg.Amount)

		result.Account, err = q.CreateAccount(ctx, CreateAccountParams{
			Owner:    arg.Owner,
			Balance:  arg.Amount,
			Currency: arg.Currency,
		})
		if err != nil {
			return err
		}

		fmt.Println(">> tx name: ", txName, "create deposit: ", arg.AccountID, "amount: ", arg.Amount)

		result.Deposit, err = q.CreateDeposit(ctx, CreateDepositParams{
			AccountID: result.Account.ID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(">> tx name: ", txName, "create audit log: ", arg.AccountID, "amount: ", arg.Amount)
		result.AuditLog, err = q.CreateLogs(ctx, CreateAuditLogParams{
			AccountID: result.Account.ID,
			Action:    "deposit",
			Timestamp: result.Deposit.CreatedAt,
		})
		if err != nil {
			return err
		}

		fmt.Println(">> tx name: ", txName, "get balance: ", arg.AccountID, "amount: ", arg.Amount)
		result.Balance, err = q.GetBalance(ctx, result.Account.ID)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

// This transaction performs the following operations atomically:
// 1. Create a new entry in the trades table representing the trade
// 2. Create a new entry in the audit_log table representing the trade
// 3. Update the balance in the accounts table
// 4. Create a new entry in the portfolio_balances table representing the trade

type TradeTxParams struct {
	AccountID int64  `json:"account_id"`
	Symbol    string `json:"symbol"`
	Amount    int32  `json:"amount"`
	Price     string `json:"price"`
	TradeType string `json:"trade_type"`
}

type TradeTxResult struct {
	Trade            Trade            `json:"trade"`
	AuditLog         AuditLog         `json:"audit_log"`
	Balance          int64            `json:"balance"`
	PortfolioBalance PortfolioBalance `json:"portfolio_balance"`
}

func (store *Store) TradeTx(ctx context.Context, arg TradeTxParams) (TradeTxResult, error) {
	var result TradeTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Trade, err = q.CreateTrade(ctx, CreateTradeParams{
			AccountID: arg.AccountID,
			Symbol:    arg.Symbol,
			Amount:    arg.Amount,
			Price:     arg.Price,
			TradeType: arg.TradeType,
		})
		if err != nil {
			return err
		}
		result.AuditLog, err = q.CreateLogs(ctx, CreateAuditLogParams{
			AccountID: arg.AccountID,
			Action:    "trade",
			Timestamp: result.Trade.CreatedAt,
		})
		if err != nil {
			return err
		}
		result.Balance, err = q.GetBalance(ctx, arg.AccountID)
		if err != nil {
			return err
		}
		result.PortfolioBalance, err = q.CreatePortfolioBalance(ctx, CreatePortfolioBalanceParams{
			AccountID: arg.AccountID,
			Symbol:    arg.Symbol,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
