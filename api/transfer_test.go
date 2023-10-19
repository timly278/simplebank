package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/timly278/simplebank/db/mock"
	db "github.com/timly278/simplebank/db/sqlc"
	"github.com/timly278/simplebank/util"
)

func TestTransferAPI(t *testing.T) {
	transferParams, transferResults := randomTransfer()

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		TransferTx(gomock.Any(), gomock.Eq(transferParams)).
		Times(1).
		Return(transferResults, nil)

	FromAccount := db.Account{
		ID:       transferParams.FromAccountID,
		Owner:    transferResults.FromAccount.Owner,
		Balance:  transferResults.FromAccount.Balance + transferParams.Amount,
		Currency: transferResults.FromAccount.Currency,
	}
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(transferParams.FromAccountID)).
		Times(1).
		Return(FromAccount, nil)

	ToAccount := db.Account{
		ID:       transferParams.ToAccountID,
		Owner:    transferResults.ToAccount.Owner,
		Balance:  transferResults.ToAccount.Balance - transferParams.Amount,
		Currency: transferResults.ToAccount.Currency,
	}
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(transferParams.ToAccountID)).
		Times(1).
		Return(ToAccount, nil)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	transferRequest := createTransferRequest{
		FromAccountID: transferParams.FromAccountID,
		ToAccountID:   transferParams.ToAccountID,
		Amount:        transferParams.Amount,
		Currency:      transferResults.FromAccount.Currency,
	}
	body, err := json.Marshal(transferRequest)
	require.NoError(t, err)

	url := "/transfers"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	requireBodyMatchTransfer(t, recorder.Body, transferResults)
}

// randomTransfer
func randomTransfer() (transferPrams db.TransferTxPrams, transferResults db.TransferTxResults) {
	account1 := randomAccount()
	account2 := randomAccount()
	account1.Currency = account2.Currency

	transferPrams = db.TransferTxPrams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	transferResults = db.TransferTxResults{
		Transfer: db.Transfer{
			ID:            util.RandomInt(1, 1000),
			FromAccountID: transferPrams.FromAccountID,
			ToAccountID:   transferPrams.ToAccountID,
			Amount:        transferPrams.Amount,
		},
		FromAccount: db.Account{
			ID:       transferPrams.FromAccountID,
			Owner:    account1.Owner,
			Balance:  account1.Balance - transferPrams.Amount,
			Currency: account1.Currency,
		},
		ToAccount: db.Account{
			ID:       transferPrams.ToAccountID,
			Owner:    account2.Owner,
			Balance:  account2.Balance + transferPrams.Amount,
			Currency: account2.Currency,
		},
		FromEntry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: transferPrams.FromAccountID,
			Amount:    -transferPrams.Amount,
		},
		ToEntry: db.Entry{
			ID:        transferResults.FromEntry.ID + 1,
			AccountID: transferPrams.ToAccountID,
			Amount:    transferPrams.Amount,
		},
	}

	return transferPrams, transferResults
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transferResults db.TransferTxResults) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.TransferTxResults

	err = json.Unmarshal(data, &gotTransfer)
	require.Equal(t, transferResults, gotTransfer)
	require.NoError(t, err)
}
