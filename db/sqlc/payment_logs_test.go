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
		Nrc:       sql.NullBool{Bool: true, Valid: true},
		BundleID:  bundle.ID,
	}
	order, _ := testQueries.CreateOrder(context.Background(), arg)
	return order
}

func createRandomPaymentLog() PaymentLog {
	order := createRandomOrder()
	arg := CreatePaymentParams{
		Payment:   rand.Float64() * 10000,
		DueDate:   time.Now().Add(time.Hour * 10),
		OrderID:   order.ID,
		Confirmed: false,
	}

	payment, _ := testQueries.CreatePayment(context.Background(), arg)
	return payment
}
func TestCreatePayment(t *testing.T) {
	order := createRandomOrder()
	arg := CreatePaymentParams{
		Payment:   rand.Float64() * 10000,
		DueDate:   time.Now().Add(time.Hour * 10),
		OrderID:   order.ID,
		Confirmed: false,
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

func TestListPaymentByConfirmation(t *testing.T) {
	for n := 0; n < 10; n++ {
		createRandomPaymentLog()
	}

	paymentLogs, err := testQueries.ListPaymentByConfirmation(context.Background(), false)
	require.NoError(t, err)
	require.NotEmpty(t, paymentLogs)
}

func TestUpdatePayment(t *testing.T) {
	paymentLog := createRandomPaymentLog()

	arg := UpdatePaymentParams{
		ID:               paymentLog.ID,
		DueDate:          paymentLog.DueDate.AddDate(0, 0, 4),
		ConfirmationDate: sql.NullTime{Time: time.Now(), Valid: true},
		OrderID:          paymentLog.OrderID,
		Confirmed:        true,
	}

	paymentLog_2, err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ID, paymentLog_2.ID)
	require.Equal(t, arg.DueDate, paymentLog_2.DueDate)
	require.WithinDuration(t, arg.ConfirmationDate.Time, paymentLog_2.ConfirmationDate.Time, time.Second)
	require.Equal(t, arg.OrderID, paymentLog_2.OrderID)
	require.Equal(t, arg.Confirmed, paymentLog_2.Confirmed)
}
