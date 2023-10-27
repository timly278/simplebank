package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/timly278/simplebank/db/mock"
	db "github.com/timly278/simplebank/db/sqlc"
	"github.com/timly278/simplebank/util"
)

type eqUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func eqParam(c db.CreateUserParams, password string) gomock.Matcher {
	return eqUserParamsMatcher{c, password}
}

// Matches returns whether x is a match.
func (e eqUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	if err := util.CheckPassword(e.password, arg.HashedPassword); err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

// String describes what the matcher matches.
func (e eqUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestCreateUser(t *testing.T) {

	password, user := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy case",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.FullName,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				userParams := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqParam(userParams, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchUser(t, user, recorder.Body)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(tc.body)
		require.NoError(t, err)

		url := "/users"
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}

}

func randomUser(t *testing.T) (string, db.User) {
	pass := util.RandomString(6)
	return pass, db.User{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}
}

func requireBodyMatchUser(t *testing.T, user db.User, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rsp userResponse
	err = json.Unmarshal(data, &rsp)
	require.NoError(t, err)

	require.NotEmpty(t, rsp)
	require.Equal(t, user.FullName, rsp.FullName)
	require.Equal(t, user.Username, rsp.Username)
	require.Equal(t, user.Email, rsp.Email)

}
