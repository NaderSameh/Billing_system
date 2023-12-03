package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/billing_system/db/mock"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	batch := db.Batch{
		ID:               1,
		Name:             "name",
		ActivationStatus: "active",
		CustomerID:       1,
		NoOfDevices:      200,
		DeliveryDate:     sql.NullTime{Time: time.Now(), Valid: true},
		WarrantyEnd:      sql.NullTime{Time: time.Now(), Valid: true},
	}
	order := db.Order{
		ID:        1,
		StartDate: time.Now(),
		EndDate:   time.Now(),
		BatchID:   batch.ID,
		BundleID:  1,
	}
	bundle := db.Bundle{
		ID:          1,
		Mrc:         1000,
		Description: "random bundle",
	}
	payment := db.PaymentLog{
		ID:      1,
		Payment: 10000,
		DueDate: time.Now(),
		OrderID: 1,
	}
	customers := db.Customer{}
	testCases := []struct {
		name          string
		body          gin.H
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK - First customer for a customer",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				// //arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(1).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "OK - Updating customers for customer",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}
				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(1).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "Missing batch name - invalid request",
			body: gin.H{
				// "batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(0).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(0).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "Internal server error - get batch",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, sql.ErrConnDone)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(0).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - create order",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, sql.ErrConnDone)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(0).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - get bundle",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(1).Return(bundle, sql.ErrConnDone)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - create payment",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(1).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, sql.ErrConnDone)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - add due",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, nil)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(1).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "Non existing batch",
			body: gin.H{
				"batch_name": batch.Name,
			},
			buildstuds: func(store *mockdb.MockStore) {
				//arg := db.AddToDueParams{Amount: bundle.Mrc * float64(batch.NoOfDevices), ID: batch.CustomerID}

				store.EXPECT().GetBatchByName(gomock.Any(), gomock.Eq(batch.Name)).Times(1).Return(batch, sql.ErrNoRows)
				store.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(bundle.ID)).Times(0).Return(bundle, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, sql.ErrConnDone)
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

			url := "/orders"

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

func TestUpdateOrder(t *testing.T) {

	randomFloat := rand.ExpFloat64()
	batch := db.Batch{
		ID:               1,
		Name:             "name",
		ActivationStatus: "active",
		CustomerID:       1,
		NoOfDevices:      100,
		DeliveryDate:     sql.NullTime{Time: time.Now(), Valid: true},
		WarrantyEnd:      sql.NullTime{Time: time.Now(), Valid: true},
	}
	bundle := db.Bundle{
		ID:          1,
		Mrc:         1000,
		Description: "random bundle",
	}
	order := db.Order{
		ID:        1,
		StartDate: time.Now().Round(time.Second),
		EndDate:   time.Now().AddDate(2, 0, 0).Round(time.Second),
		BatchID:   batch.ID,
		BundleID:  10,
		Nrc:       sql.NullFloat64{Float64: randomFloat, Valid: true},
	}
	payment := db.PaymentLog{
		ID:      1,
		Payment: 10000,
		DueDate: time.Now(),
		OrderID: 1,
	}
	customers := db.Customer{}

	testCases := []struct {
		name          string
		body          gin.H
		OrderID       int64
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK - with NRC",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "OK - without NRC",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        sql.NullFloat64{Float64: randomFloat, Valid: false},
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				// arg := db.UpdateOrdersParams{
				// 	ID:        order.ID,
				// 	Nrc:       sql.NullFloat64{Float64: randomFloat, Valid: false},
				// 	BundleID:  order.BundleID,
				// 	StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
				// 	EndDate:   order.EndDate,
				// }
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "bad request - json",
			body: gin.H{
				"end_date": order.EndDate,
				// "bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(0).Return(order, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "bad request - uri",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			// OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(0).Return(order, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "Order not found",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, sql.ErrNoRows)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(0).Return(bundle, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name: "Internal server error - updating order",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, sql.ErrConnDone)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(0).Return(bundle, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "Internal server error - get bundle",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, sql.ErrConnDone)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(0).Return(batch, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - get batch",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, sql.ErrConnDone)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				// store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
				// store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				// store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - create payment",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, sql.ErrConnDone)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)
				// store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				// store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error - update customer",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, sql.ErrConnDone)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
				// store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "Internal server error - create payment NRC",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, sql.ErrConnDone)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(customers, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "Internal server error - update customers NRC",
			body: gin.H{
				"end_date":   order.EndDate,
				"bundle_id":  order.BundleID,
				"nrc":        order.Nrc,
				"start_date": sql.NullTime{Time: order.StartDate, Valid: true},
			},
			OrderID: order.ID,
			buildstuds: func(store *mockdb.MockStore) {
				/*arg := db.UpdateOrdersParams{
					ID:        order.ID,
					Nrc:       order.Nrc,
					BundleID:  order.BundleID,
					StartDate: sql.NullTime{Time: order.StartDate, Valid: true},
					EndDate:   order.EndDate,
				} */
				store.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Times(1).Return(order, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Eq(order.BundleID)).Times(1).Return(bundle, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(batch, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(payment, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(customers, sql.ErrConnDone)

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

			url := fmt.Sprintf("/orders/%d", tc.OrderID)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
