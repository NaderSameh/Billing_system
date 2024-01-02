package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
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

func TestUpdatePaymentLog(t *testing.T) {

	payment := db.PaymentLog{
		ID:         rand.Int63(),
		Payment:    math.Round(rand.ExpFloat64() * 100),
		DueDate:    time.Now().Round(time.Second),
		OrderID:    rand.Int63(),
		Confirmed:  false,
		CustomerID: rand.Int63(),
	}

	testCases := []struct {
		name          string
		PaymentID     int64
		body          gin.H
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:               payment.ID,
					ConfirmationDate: sql.NullTime{Time: time.Now().Round(time.Second), Valid: true},
					DueDate:          sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					Confirmed:        true,
				}

				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().AddToPaid(gomock.Any(), gomock.AssignableToTypeOf(db.AddToPaidParams{Amount: payment.Payment, ID: payment.CustomerID})).Times(1).Return(db.Customer{}, nil)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:      "Invalid params - missing confirmation",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				// "confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "Invalid params - no path ID",
			// PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name:      "OK with false confirmation",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": false,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:        payment.ID,
					DueDate:   sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					Confirmed: false,
				}
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "OK with true confirmation on previously confirmed",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:               payment.ID,
					DueDate:          sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					ConfirmationDate: sql.NullTime{Time: time.Now().Round(time.Second), Valid: true},
					Confirmed:        true,
				}
				payment.Confirmed = true
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:      "Non existing payment",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:        payment.ID,
					DueDate:   sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					Confirmed: payment.Confirmed,
				}
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, sql.ErrNoRows)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name:      "Internal server error - get payment",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:        payment.ID,
					DueDate:   sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					Confirmed: payment.Confirmed,
				}
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, sql.ErrConnDone)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal server error - add to paid",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				payment.Confirmed = false
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().AddToPaid(gomock.Any(), gomock.Eq(db.AddToPaidParams{Amount: payment.Payment, ID: payment.CustomerID})).Times(1).Return(db.Customer{}, sql.ErrConnDone)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.Any()).Times(0).Return(payment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name:      "Internal server error - update payment",
			PaymentID: payment.ID,
			body: gin.H{
				"due_date": gin.H{
					"time":  payment.DueDate.Round(time.Second),
					"valid": true,
				},
				"confirmed": true,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentParams{
					ID:               payment.ID,
					ConfirmationDate: sql.NullTime{Time: time.Now().Round(time.Second), Valid: true},
					DueDate:          sql.NullTime{Time: payment.DueDate.Round(time.Second), Valid: true},
					Confirmed:        true,
				}
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().AddToPaid(gomock.Any(), gomock.Eq(db.AddToPaidParams{Amount: payment.Payment, ID: payment.CustomerID})).Times(1).Return(db.Customer{}, nil)
				store.EXPECT().UpdatePayment(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payment, sql.ErrConnDone)
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
			data, err := json.Marshal(tc.body)
			url := fmt.Sprintf("/payments_logs/%d", tc.PaymentID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestDeletePayment(t *testing.T) {

	payment := db.PaymentLog{
		ID:         rand.Int63(),
		Payment:    math.Round(rand.ExpFloat64() * 100),
		DueDate:    time.Now().Round(time.Second),
		OrderID:    rand.Int63(),
		Confirmed:  true,
		CustomerID: rand.Int63(),
	}
	testCases := []struct {
		name          string
		PaymentID     int64
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{

			name:      "OK",
			PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.AddToDueParams{
					Amount: -payment.Payment,
					ID:     payment.CustomerID,
				}
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(db.Customer{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "invalid request - No ID",
			// PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(0).Return(payment, nil)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "No payment record",
			PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, sql.ErrNoRows)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal server err - get payment",
			PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, sql.ErrConnDone)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal server err - delete payment",
			PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal server err - add to due",
			PaymentID: payment.ID,
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().GetPaymentForUpdate(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(payment, nil)
				store.EXPECT().DeletePayment(gomock.Any(), gomock.Eq(payment.ID)).Times(1).Return(nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(db.Customer{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/payments_logs/%d", tc.PaymentID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestListPayments(t *testing.T) {

	payments := make([]db.PaymentLog, 10)
	customerID := rand.Int63()
	payments[0] = db.PaymentLog{
		ID:        rand.Int63(),
		Payment:   math.Round(rand.ExpFloat64() * 100),
		DueDate:   time.Now().Round(time.Second),
		OrderID:   rand.Int63(),
		Confirmed: true,
	}

	type Query struct {
		Confirmed  *bool
		CustomerID *int64
		PageID     int32
		PageSize   int32
	}

	testCases := []struct {
		name          string
		query         Query
		PaymentID     int64
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK - true confirmation",
			query: Query{
				Confirmed:  boolPtr(true),
				CustomerID: &customerID,
				PageID:     1,
				PageSize:   5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
					CustomerID: sql.NullInt64{
						Int64: customerID,
						Valid: true,
					},
				}

				arg2 := db.ListAllPaymentsCountParams{
					Confirmed:  arg.Confirmed,
					CustomerID: arg.CustomerID,
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Eq(arg2)).Times(1).Return(int64(10), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK - false confirmation",
			query: Query{
				Confirmed:  boolPtr(false),
				CustomerID: &customerID,
				PageID:     1,
				PageSize:   5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
					CustomerID: sql.NullInt64{
						Int64: customerID,
						Valid: true,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Any()).Times(1).Return(int64(11), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK - no confirmation",
			query: Query{
				Confirmed:  nil,
				CustomerID: &customerID,
				PageID:     1,
				PageSize:   5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  false,
						Valid: false,
					},
					CustomerID: sql.NullInt64{
						Int64: customerID,
						Valid: true,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Any()).Times(1).Return(int64(10), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "OK - no customer ID",
			query: Query{
				Confirmed: boolPtr(true),
				// CustomerID: &customerID,
				PageID:   1,
				PageSize: 5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
					CustomerID: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Any()).Times(1).Return(int64(10), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "OK - no optional params",
			query: Query{
				// Confirmed: boolPtr(true),
				// CustomerID: &customerID,
				PageID:   1,
				PageSize: 5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  false,
						Valid: false,
					},
					CustomerID: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Any()).Times(1).Return(int64(10), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal server error - count",
			query: Query{
				// Confirmed: boolPtr(true),
				// CustomerID: &customerID,
				PageID:   1,
				PageSize: 5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  false,
						Valid: false,
					},
					CustomerID: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, nil)
				store.EXPECT().ListAllPaymentsCount(gomock.Any(), gomock.Any()).Times(1).Return(int64(10), sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid params",
			query: Query{
				Confirmed:  boolPtr(true),
				CustomerID: &customerID,
				PageID:     0,
				PageSize:   5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(0).Return(payments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error",
			query: Query{
				Confirmed:  boolPtr(true),
				CustomerID: &customerID,
				PageID:     1,
				PageSize:   5,
			},
			buildstuds: func(store *mockdb.MockStore) {
				arg := db.ListPaymentsParams{
					Confirmed: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
					CustomerID: sql.NullInt64{
						Int64: customerID,
						Valid: true,
					},
				}
				store.EXPECT().ListPayments(gomock.Any(), gomock.AssignableToTypeOf(arg)).Times(1).Return(payments, sql.ErrConnDone)
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

			url := "/payments_logs"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			values := request.URL.Query()
			if tc.query.Confirmed != nil {
				values.Add("confirmed", fmt.Sprintf("%t", *tc.query.Confirmed))
			}
			if tc.query.CustomerID != nil {
				values.Add("customer_id", fmt.Sprintf("%d", *tc.query.CustomerID))
			}
			values.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			values.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = values.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
