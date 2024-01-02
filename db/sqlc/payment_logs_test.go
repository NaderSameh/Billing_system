package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomOrder() Order {
	batch := createRandomBatch()
	bundle := createRandomBundle()
	arg := CreateOrderParams{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 50),
		BatchID:   batch.ID,
		Nrc:       sql.NullFloat64{Float64: 10000.50, Valid: true},
		BundleID:  bundle.ID,
	}
	order, _ := testQueries.CreateOrder(context.Background(), arg)
	return order
}

func createRandomPaymentLog() PaymentLog {
	order := createRandomOrder()
	customer := createRandomCustomer()
	arg := CreatePaymentParams{
		Payment:    rand.Float64() * 10000,
		DueDate:    time.Now().Add(time.Hour * 10),
		OrderID:    order.ID,
		Confirmed:  false,
		CustomerID: customer.ID,
	}

	payment, _ := testQueries.CreatePayment(context.Background(), arg)
	return payment
}
func TestCreatePayment(t *testing.T) {
	order := createRandomOrder()
	customer := createRandomCustomer()
	arg := CreatePaymentParams{
		Payment:    rand.Float64() * 10000,
		DueDate:    time.Now().Add(time.Hour * 10),
		OrderID:    order.ID,
		Confirmed:  false,
		CustomerID: customer.ID,
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, payment.Payment, arg.Payment)
	require.Equal(t, payment.OrderID, arg.OrderID)
	require.Equal(t, payment.Confirmed, arg.Confirmed)
	require.WithinDuration(t, payment.DueDate, arg.DueDate, time.Second)
}

func TestDeletePayment(t *testing.T) {
	paymentLog := createRandomPaymentLog()
	err := testQueries.DeletePayment(context.Background(), paymentLog.ID)
	require.NoError(t, err)
}

func TestGetPaymentForUpdate(t *testing.T) {
	paymentLog := createRandomPaymentLog()
	payment_log2, err := testQueries.GetPaymentForUpdate(context.Background(), paymentLog.ID)
	require.NoError(t, err)
	require.Equal(t, paymentLog, payment_log2)
}

func TestListPayments(t *testing.T) {
	for n := 0; n < 10; n++ {
		createRandomPaymentLog()
	}
	arg := ListPaymentsParams{
		Confirmed:  sql.NullBool{Bool: true, Valid: false},
		CustomerID: sql.NullInt64{Int64: rand.Int63(), Valid: false},
		Limit:      5,
		Offset:     0,
	}
	paymentLogs, err := testQueries.ListPayments(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paymentLogs)

	arg2 := ListPaymentsParams{
		Confirmed:  sql.NullBool{Bool: true, Valid: true},
		CustomerID: sql.NullInt64{Int64: rand.Int63(), Valid: true},
	}
	paymentLogs, err = testQueries.ListPayments(context.Background(), arg2)
	require.NoError(t, err)

}

func TestUpdatePayment(t *testing.T) {
	paymentLog := createRandomPaymentLog()

	arg := UpdatePaymentParams{
		ID:      paymentLog.ID,
		DueDate: sql.NullTime{Time: paymentLog.DueDate.AddDate(0, 0, 4), Valid: true},

		Confirmed: true,
	}

	paymentLog_2, err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ID, paymentLog_2.ID)
	require.Equal(t, arg.DueDate.Time, paymentLog_2.DueDate)
	require.Equal(t, arg.Confirmed, paymentLog_2.Confirmed)
	arg.ID = rand.Int63() + 200
	paymentLog_2, err = testQueries.UpdatePayment(context.Background(), arg)
	require.Equal(t, err, sql.ErrNoRows)
}
