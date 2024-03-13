package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Printf("Before\nFrom account balance: %v\nTo account balance: %v\n", fromAccount.Balance, toAccount.Balance)

	parallel := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < parallel; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < parallel; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		from := result.FromAccount
		require.NotEmpty(t, from)
		require.Equal(t, fromAccount.ID, from.ID)

		to := result.ToAccount
		require.NotEmpty(t, to)
		require.Equal(t, toAccount.ID, to.ID)

		// check balances
		fmt.Printf("After tx %v >>\nFrom account balance: %v\nTo account balance: %v\n", i, from.Balance, to.Balance)
		diff1 := fromAccount.Balance - from.Balance
		diff2 := to.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= parallel)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	afterFromAccount, err := testStore.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	afterToAccount, err := testStore.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Printf("After\nFrom account balance: %v\nTo account balance: %v\n", afterFromAccount.Balance, afterToAccount.Balance)

	require.Equal(t, afterFromAccount.Balance, fromAccount.Balance-amount*int64(parallel))
	require.Equal(t, afterToAccount.Balance, toAccount.Balance+amount*int64(parallel))
}
