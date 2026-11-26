

package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount := randomAccount()
	toAccount := randomAccount()

	// Ensure both accounts use the same currency.
	fromAccount.Currency = "USD"
	toAccount.Currency = "USD"

	amount := util.RandomInt(1, 1000)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccount.ID).
					Times(1).
					Return(toAccount, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidFromAccount",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidToAccount",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccount.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "EUR",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				wrongCurrencyAccount := toAccount
				wrongCurrencyAccount.Currency = "EUR"

				store.EXPECT().
					GetAccount(gomock.Any(), toAccount.ID).
					Times(1).
					Return(wrongCurrencyAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountDoesNotBelongToUser",
			body: map[string]interface{}{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":           amount,
				"currency":         "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(
				http.MethodPost,
				"/transfer",
				strings.NewReader(string(body)),
			)
			require.NoError(t, err)

			username := fromAccount.Owner
			if tc.name == "FromAccountDoesNotBelongToUser" {
				username = toAccount.Owner
			}

			addAuthorization(
				t,
				request,
				server.tokenMaker,
				authorizationTypeBearer,
				username,
				time.Minute,
			)

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
