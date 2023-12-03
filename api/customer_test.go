package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/billing_system/db/mock"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestListCharges(t *testing.T) {
	customer := db.Customer{
		ID:       1,
		Customer: "random",
	}
	customers := make([]db.Customer, 10)
	testCases := []struct {
		name          string
		CustomerName  string
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			CustomerName: customer.Customer,
			buildstuds: func(store *mockdb.MockStore) {
				arg := sql.NullString{String: customer.Customer, Valid: true}

				store.EXPECT().ListAllCharges(gomock.Any(), gomock.Eq(arg)).Times(1).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "List all bundles - no customer name",
			// CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				arg := sql.NullString{String: "", Valid: false}
				store.EXPECT().ListAllCharges(gomock.Any(), gomock.Eq(arg)).Times(1).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:         "Internal server error - list charges",
			CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				arg := sql.NullString{String: customer.Customer, Valid: true}
				store.EXPECT().ListAllCharges(gomock.Any(), gomock.Eq(arg)).Times(1).Return(customers, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - list charges - no name",
			// CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				arg := sql.NullString{String: "", Valid: false}
				store.EXPECT().ListAllCharges(gomock.Any(), gomock.Eq(arg)).Times(1).Return(customers, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstuds(store)

			server := newTestServer(t, store)
			url := "/charges"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("customer_name", fmt.Sprintf("%s", tc.CustomerName))

			request.URL.RawQuery = q.Encode()
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
