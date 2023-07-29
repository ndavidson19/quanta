package db

import (
	"context"
	"testing"
)

func TestCreateTx(t *testing.T) {
	store := NewStore(testDB)

	account := CreateAccountParams{
		Owner: "test",
		Balance: 0,
		Currency: "USD",
	}

	// run n concurrent transactions
	n := 5
	ammount := int64(10)

	errs := make(chan error)
	results := make(chan CreateTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.CreateTx(context.Background(), CreateTxParams{
				AccountID: account.ID,
				Amount: ammount,
			})

			errs <- err
			results <- result

		}()
	}

	// check results
	for i := 0, i < n;, i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check account
		require.NotEmpty(t, result.Account)
		require.Equal(t, account.Owner, result.Account.Owner)
		require.Equal(t, account.Balance + ammount, result.Account.Balance)
		require.Equal(t, account.Currency, result.Account.Currency)

		// check deposit
		require.NotEmpty(t, result.Deposit)
		require.Equal(t, account.ID, result.Deposit.AccountID)
		require.Equal(t, ammount, result.Deposit.Amount)

		// check audit log
		require.NotEmpty(t, result.AuditLog)
		require.Equal(t, account.ID, result.AuditLog.AccountID)
		require.Equal(t, "deposit", result.AuditLog.Action)
		require.NotEmpty(t, result.AuditLog.Timestamp)

		// check balance
		require.NotEmpty(t, result.Balance)
		require.Equal(t, account.Balance + ammount, result.Balance)
			
		_, err = store.GetAccount(context.Background(), result.Account.ID)
		require.NoError(t, err)

		_, err = store.GetDeposit(context.Background(), result.Deposit.ID)
		require.NoError(t, err)

		_, err = store.GetAuditLog(context.Background(), result.AuditLog.ID)
		require.NoError(t, err)

	}
}

