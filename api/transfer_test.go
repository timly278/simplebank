package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/timly278/simplebank/db/mock"
	db "github.com/timly278/simplebank/db/sqlc"
	"github.com/timly278/simplebank/token"
	"github.com/timly278/simplebank/util"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)
	_, user1 := randomUser(t)
	_, user2 := randomUser(t)
	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)

	account1.Currency = util.USD
	account2.Currency = util.USD

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy Case",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				transferParams := db.TransferTxPrams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(transferParams)).Times(1)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		ctrl := gomock.NewController(t)
		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		// start test server and send request
		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(tc.body)
		require.NoError(t, err)

		url := "/transfers"
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		require.NoError(t, err)

		tc.setupAuth(t, request, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}
