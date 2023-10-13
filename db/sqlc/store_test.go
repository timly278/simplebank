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
	n := 5
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxPrams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			ech <-err
			rch <-result
		}()

	}

	for i := 0; i < n; i++ {
		err := <- ech
		require.NoError(t, err)
		
		result := <- rch

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

		// check to 
		entry = result.ToEntry
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, amount, entry.Amount)
		require.Equal(t, toAccount.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
	
		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		// TODO: update accounts' balance
	}

}
