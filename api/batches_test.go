package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/billing_system/db/mock"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/util"

	"github.com/stretchr/testify/require"
)

func createRandomBatch() db.Batch {

	batch := db.Batch{
		ID:               1,
		Name:             "random",
		ActivationStatus: "active",
		CustomerID:       1,
		NoOfDevices:      100,
		DeliveryDate:     sql.NullTime{Time: time.Now().Round(time.Second), Valid: true},
		WarrantyEnd:      sql.NullTime{Time: time.Now().Round(time.Second), Valid: true},
	}
	return batch
}

func TestCreateBatch(t *testing.T) {

	batch := createRandomBatch()
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
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     batch.NoOfDevices,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Non existing customer",
			body: gin.H{
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     batch.NoOfDevices,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().CreateCustomer(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, nil)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid param",
			body: gin.H{
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     0,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(0).Return(customer, nil)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     batch.NoOfDevices,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - create batch",
			body: gin.H{
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     batch.NoOfDevices,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batch, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - create customer",
			body: gin.H{
				"name":              batch.Name,
				"activation_status": batch.ActivationStatus,
				"customer_name":     customer.Customer,
				"no_of_devices":     batch.NoOfDevices,
				"delivery_date":     time.Now().Round(time.Second),
				"warranty_end":      time.Now().Round(time.Second),
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.CreateBatchParams{
					Name:             batch.Name,
					ActivationStatus: batch.ActivationStatus,
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices,
					DeliveryDate:     batch.DeliveryDate,
					WarrantyEnd:      batch.WarrantyEnd,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().CreateCustomer(gomock.Any(), gomock.Eq(customer.Customer)).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().CreateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
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

			url := "/batches"

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
func TestDeleteBatch(t *testing.T) {
	batch := createRandomBatch()
	type Query struct {
		batch_id int64
	}

	testCases := []struct {
		name          string
		query         Query
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				batch_id: batch.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(batch, nil)
				store.EXPECT().DeleteBatch(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "batch not found",
			query: Query{
				batch_id: batch.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(batch, sql.ErrNoRows)
				store.EXPECT().DeleteBatch(gomock.Any(), gomock.Eq(batch.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name: "invalid id",
			query: Query{
				batch_id: 0,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Eq(batch.ID)).Times(0).Return(batch, nil)
				store.EXPECT().DeleteBatch(gomock.Any(), gomock.Eq(batch.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error - get",
			query: Query{
				batch_id: batch.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(batch, sql.ErrConnDone)
				store.EXPECT().DeleteBatch(gomock.Any(), gomock.Eq(batch.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "internal server error - delete",
			query: Query{
				batch_id: batch.ID,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(batch, nil)
				store.EXPECT().DeleteBatch(gomock.Any(), gomock.Eq(batch.ID)).Times(1).Return(sql.ErrConnDone)
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
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/batches/%d", tc.query.batch_id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestListBatches(t *testing.T) {
	n := 10
	batches := make([]db.Batch, n)
	for i := 0; i < n; i++ {
		batches[i] = createRandomBatch()
	}

	customer := db.Customer{
		ID:       1,
		Customer: "random customer",
	}
	type Query struct {
		CustomerName string
		page_id      int32
		page_size    int32
	}

	testCases := []struct {
		name          string
		query         Query
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				CustomerName: "random",
				page_id:      1,
				page_size:    5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:      5,
					Offset:     0,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batches, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "Invalid param",
			query: Query{
				CustomerName: "random",
				page_id:      0,
				page_size:    5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:      5,
					Offset:     0,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(0).Return(customer, nil)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batches, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			query: Query{
				CustomerName: "random",
				page_id:      1,
				page_size:    5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:      5,
					Offset:     0,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batches, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - listing",
			query: Query{
				CustomerName: "random",
				page_id:      1,
				page_size:    5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:      5,
					Offset:     0,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batches, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "no customer name param",
			query: Query{
				// CustomerName: "random",
				page_id:   1,
				page_size: 5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					// CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:  5,
					Offset: 0,
				}
				// store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batches, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "non existing customer name",
			query: Query{
				CustomerName: "random",
				page_id:      1,
				page_size:    5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListAllBatchesParams{
					CustomerID: sql.NullInt64{Valid: true, Int64: 1},
					Limit:      5,
					Offset:     0,
				}
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().ListAllBatches(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batches, nil)
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
			recorder := httptest.NewRecorder()

			url := "/batches"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("customer_name", fmt.Sprintf("%s", tc.query.CustomerName))
			q.Add("page_id", fmt.Sprintf("%d", tc.query.page_id))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.page_size))
			request.URL.RawQuery = q.Encode()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestUpdateBatches(t *testing.T) {

	batch := createRandomBatch()
	customer := db.Customer{
		ID:       1,
		Customer: "random customer",
	}
	testCases := []struct {
		name          string
		BatchID       int64
		body          gin.H
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "Invalid param URI",
			BatchID: -10,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(0).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(0).Return(customer, nil)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name:    "Invalid param - body",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     10,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(0).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(0).Return(customer, nil)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:    "Internal server error - get",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrConnDone)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "Internal server error - update",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, nil)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(batch, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:    "Non existing customer name",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:    "Invalid param",
			BatchID: batch.ID,
			body: gin.H{
				"customer_name":     customer.Customer,
				"activation_status": "suspended",
				"no_of_devices":     batch.NoOfDevices * 10,
				"delivery_date":     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				"warranty_end":      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdateBatchParams{
					ID:               batch.ID,
					ActivationStatus: "suspended",
					WarrantyEnd:      sql.NullTime{Time: batch.WarrantyEnd.Time.AddDate(1, 1, 0), Valid: true},
					CustomerID:       batch.CustomerID,
					NoOfDevices:      batch.NoOfDevices * 10,
					DeliveryDate:     sql.NullTime{Time: batch.DeliveryDate.Time.AddDate(0, 1, 0), Valid: true},
				}
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().GetCustomerID(gomock.Any(), gomock.Any()).Times(1).Return(customer, sql.ErrNoRows)
				store.EXPECT().UpdateBatch(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(batch, nil)
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
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			url := fmt.Sprintf("/batches/%d", tc.BatchID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
