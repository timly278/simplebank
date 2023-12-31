package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// Write test, create 5 or 10 transactions of 2 accounts by using goroutines

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	rch := make(chan TransferTxResults)
	ech := make(chan error)
	// namechan := make(chan string)
	n := 5
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxPrams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			ech <- err
			rch <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-ech
		require.NoError(t, err)

		result := <-rch
		// txName := <-namechan
		// check transfer
		transfer := result.Transfer
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entry
		entry := result.FromEntry
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, -amount, entry.Amount)
		require.Equal(t, fromAccount.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		// check to entry
		entry = result.ToEntry
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, amount, entry.Amount)
		require.Equal(t, toAccount.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		// Check accounts' ID
		acc1 := result.FromAccount
		require.NotEmpty(t, acc1)
		require.Equal(t, fromAccount.ID, acc1.ID)

		acc2 := result.ToAccount
		require.NotEmpty(t, acc2)
		require.Equal(t, toAccount.ID, acc2.ID)

		// Check accounts' balance
		diff1 := fromAccount.Balance - acc1.Balance
		diff2 := acc2.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, ... n * amount

		// diff1/amount = 1, 2, 3, 4, 5
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// check the final update balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferDeadlock(t *testing.T) {

	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	ech := make(chan error)
	n := 10
	amount := int64(10)
	for i := 0; i < n; i++ {
		fromID := fromAccount.ID
		toID := toAccount.ID

		if i % 2 == 1 {
			fromID = toAccount.ID
			toID = fromAccount.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxPrams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})
			ech <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <- ech
		require.NoError(t, err)
	}
	// check the final update balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance, updatedAccount2.Balance)

}
