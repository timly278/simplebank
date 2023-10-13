package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	amount := int64(10)
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        amount,
	})
	require.NoError(t, err)
	require.NotZero(t, transfer.ID)
	require.Equal(t, acc1.ID, transfer.FromAccountID)
	require.Equal(t, acc2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	tr, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, tr)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestGetListTransfers(t *testing.T) {
	// TODO:
}

