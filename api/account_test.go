package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder, account)
			},
		},
		{
			name:      "NotFound",
			accountID: 999999,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(999999))).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:      util.RandomInt(1, 1000),
		Owner:   util.RandomOwner(),
		Balance: util.RandomInt(1, 1000),
	}
}

func requireBodyMatchAccount(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
	var gotAccount db.Account

	err := json.Unmarshal(recorder.Body.Bytes(), &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
