package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/timly278/simplebank/util"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	entryPram := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInt(1, 1000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), entryPram)

	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)
	require.Equal(t, entryPram.AccountID, entry.AccountID)
	require.Equal(t, entryPram.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.Empty(t, entry2)

}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.ID, entry2.ID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestGetListEntry(t *testing.T) {

	var randomParam1 CreateEntryParams
	var randomParam2 CreateEntryParams

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	randomParam2.AccountID = int64(account1.ID)
	randomParam1.AccountID = int64(account2.ID)

	for i := 0; i < 5; i++ {
		randomParam1.Amount = util.RandomInt(1, 1000)
		randomParam2.Amount = util.RandomInt(1, 1000)

		_, err := testQueries.CreateEntry(context.Background(), randomParam1)
		require.NoError(t, err)

		_, err = testQueries.CreateEntry(context.Background(), randomParam2)
		require.NoError(t, err)
	}

	listParam := ListEntriesParams{
		AccountID: randomParam2.AccountID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), listParam)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, v := range entries {
		require.NotEmpty(t, v)
	}
}
