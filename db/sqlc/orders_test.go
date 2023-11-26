package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	batch := createRandomBatch()
	bundle := createRandomBundle()
	arg := CreateOrderParams{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 50),
		BatchID:   batch.ID,
		Nrc:       sql.NullBool{Bool: true, Valid: true},
		BundleID:  bundle.ID,
	}
	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.WithinDuration(t, order.StartDate, arg.StartDate, time.Second)
	require.WithinDuration(t, order.EndDate, arg.EndDate, time.Second)
	require.Equal(t, order.BatchID, arg.BatchID)
	require.Equal(t, order.Nrc, arg.Nrc)
	require.Equal(t, order.BundleID, arg.BundleID)

}

func TestDeleteOrder(t *testing.T) {
	batch := createRandomBatch()
	bundle := createRandomBundle()
	arg := CreateOrderParams{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 50),
		BatchID:   batch.ID,
		Nrc:       sql.NullBool{Bool: true, Valid: true},
		BundleID:  bundle.ID,
	}
	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	err = testQueries.DeleteOrder(context.Background(), order.ID)
	require.NoError(t, err)

}

func TestListOrdersByBatchID(t *testing.T) {

	batch := createRandomBatch()
	bundle := createRandomBundle()
	for n := 0; n < 10; n++ {
		arg := CreateOrderParams{
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 50),
			BatchID:   batch.ID,
			Nrc:       sql.NullBool{Bool: true, Valid: true},
			BundleID:  bundle.ID,
		}
		_, err := testQueries.CreateOrder(context.Background(), arg)
		require.NoError(t, err)
	}

	orders, err := testQueries.ListOrdersByBatchID(context.Background(), batch.ID)
	require.NoError(t, err)
	require.Len(t, orders, 10)
}

func TestListOrdersByBundleID(t *testing.T) {

	batch := createRandomBatch()
	bundle := createRandomBundle()
	for n := 0; n < 10; n++ {
		arg := CreateOrderParams{
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 50),
			BatchID:   batch.ID,
			Nrc:       sql.NullBool{Bool: true, Valid: true},
			BundleID:  bundle.ID,
		}
		_, err := testQueries.CreateOrder(context.Background(), arg)
		require.NoError(t, err)
	}

	orders, err := testQueries.ListOrdersByBundleID(context.Background(), bundle.ID)
	require.NoError(t, err)
	require.Len(t, orders, 10)
}

func TestUpdateOrders(t *testing.T) {

	batch := createRandomBatch()
	bundle := createRandomBundle()

	arg := CreateOrderParams{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 50),
		BatchID:   batch.ID,
		Nrc:       sql.NullBool{Bool: true, Valid: true},
		BundleID:  bundle.ID,
	}
	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)

	arg2 := UpdateOrdersParams{
		ID:        order.ID,
		Nrc:       sql.NullBool{Bool: false, Valid: false},
		BundleID:  bundle.ID,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 100),
	}

	orderUpdated, err := testQueries.UpdateOrders(context.Background(), arg2)
	require.NoError(t, err)
	require.Equal(t, orderUpdated.Nrc, arg2.Nrc)
	require.Equal(t, orderUpdated.BundleID, arg2.BundleID)
	require.WithinDuration(t, orderUpdated.StartDate, arg2.StartDate, time.Second)
	require.WithinDuration(t, orderUpdated.EndDate, arg2.EndDate, time.Second)
}
