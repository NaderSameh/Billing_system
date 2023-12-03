package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/billing_system/db/mock"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func TestCreateBundles(t *testing.T) {
	bundle := db.Bundle{
		Mrc:         1000,
		Description: "Simplest bundle of 1000$",
	}
	testCases := []struct {
		name          string
		body          gin.H
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"mrc":         bundle.Mrc,
				"description": bundle.Description,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBundleParams{
					Mrc:         bundle.Mrc,
					Description: bundle.Description,
				}
				store.EXPECT().CreateBundle(gomock.Any(), gomock.Eq(arg)).Times(1).Return(bundle, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid params",
			body: gin.H{
				"mrc":         -bundle.Mrc,
				"description": bundle.Description,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBundleParams{
					Mrc:         bundle.Mrc,
					Description: bundle.Description,
				}
				store.EXPECT().CreateBundle(gomock.Any(), gomock.Eq(arg)).Times(0).Return(bundle, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"mrc":         bundle.Mrc,
				"description": bundle.Description,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBundleParams{
					Mrc:         bundle.Mrc,
					Description: bundle.Description,
				}
				store.EXPECT().CreateBundle(gomock.Any(), gomock.Eq(arg)).Times(1).Return(bundle, sql.ErrConnDone)
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

			url := "/bundles"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestAssignBundle(t *testing.T) {
	bundle := db.Bundle{
		ID:          1,
		Mrc:         1000,
		Description: "Simplest bundle of 1000$",
	}
	customer := db.Customer{
		ID:       1,
		Customer: util.GenerateRandomString(10),
	}
	testCases := []struct {
		name          string
		body          gin.H
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"customer_name": customer.Customer,
				"bundle_id":     bundle.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: customer.ID}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, nil)
				store.EXPECT().AddBundleToCustomer(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "invalid param",
			body: gin.H{
				// "customer_name": customer.Customer,
				"bundle_id": bundle.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: customer.ID}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq(customer.Customer)).Times(0).Return(customer, nil)
				store.EXPECT().AddBundleToCustomer(gomock.Any(), gomock.Eq(arg)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error - get customer ID",
			body: gin.H{
				"customer_name": customer.Customer,
				"bundle_id":     bundle.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: customer.ID}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().AddBundleToCustomer(gomock.Any(), gomock.Eq(arg)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "internal server error - assign bundle",
			body: gin.H{
				"customer_name": customer.Customer,
				"bundle_id":     bundle.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: customer.ID}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, nil)
				store.EXPECT().AddBundleToCustomer(gomock.Any(), gomock.Eq(arg)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "internal server error - assign bundle",
			body: gin.H{
				"customer_name": customer.Customer,
				"bundle_id":     bundle.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				// arg := db.AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: customer.ID}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, sql.ErrNoRows)
				// store.EXPECT().AddBundleToCustomer(gomock.Any(), gomock.Eq(arg)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := "/bundles/assign"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestGetBundles(t *testing.T) {
	customer := db.Customer{
		ID:       1,
		Customer: "random",
	}
	bundles := make([]db.Bundle, 10)
	testCases := []struct {
		name          string
		CustomerName  string
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq("random")).Times(1).Return(customer, nil)
				store.EXPECT().ListBundlesByCustomerID(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(bundles, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "List all bundles - no customer name",
			// CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().ListAllBundles(gomock.Any()).Times(1).Return(bundles, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:         "Internal server error",
			CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq("random")).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().ListBundlesByCustomerID(gomock.Any(), gomock.Eq(customer.ID)).Times(0).Return(bundles, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "Internal server error - get",
			CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq("random")).Times(1).Return(customer, nil)
				store.EXPECT().ListBundlesByCustomerID(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(bundles, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - get - no name",
			// CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				// store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq("random")).Times(1).Return(customer, nil)
				store.EXPECT().ListAllBundles(gomock.Any()).Times(1).Return(bundles, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:         "Non existing customer",
			CustomerName: "random",
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Eq("random")).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().ListBundlesByCustomerID(gomock.Any(), gomock.Eq(customer.ID)).Times(0).Return(bundles, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
			url := "/bundles"
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

func TestDeleteBundle(t *testing.T) {
	bundle := db.Bundle{
		ID:          10,
		Mrc:         200,
		Description: "random",
	}
	testCases := []struct {
		name          string
		BundleID      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			BundleID: bundle.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetBundleByID(gomock.Any(), bundle.ID).Times(1)
				store.EXPECT().DeleteBundle(gomock.Any(), bundle.ID).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "not found",
			BundleID: bundle.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetBundleByID(gomock.Any(), bundle.ID).Times(1).Return(db.Bundle{}, sql.ErrNoRows)
				store.EXPECT().DeleteBundle(gomock.Any(), bundle.ID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "not found",
			// BundleID: bundle.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetBundleByID(gomock.Any(), bundle.ID).Times(0).Return(db.Bundle{}, sql.ErrNoRows)
				store.EXPECT().DeleteBundle(gomock.Any(), bundle.ID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "internal server error get",
			BundleID: bundle.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetBundleByID(gomock.Any(), bundle.ID).Times(1).Return(db.Bundle{}, sql.ErrConnDone)
				store.EXPECT().DeleteBundle(gomock.Any(), bundle.ID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "internal server error delete",
			BundleID: bundle.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetBundleByID(gomock.Any(), bundle.ID).Times(1).Return(db.Bundle{}, nil)
				store.EXPECT().DeleteBundle(gomock.Any(), bundle.ID).Times(1).Return(sql.ErrConnDone)
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
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/bundles/%d", tc.BundleID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
