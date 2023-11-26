package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func createRandomBundle() Bundle {
	arg := CreateBundleParams{
		Mrc:         rand.Int31n(100),
		Description: util.GenerateRandomString(20),
	}
	bundle, _ := testQueries.CreateBundle(context.Background(), arg)
	return bundle
}

func createRandomCustomer() Customer {
	customer, _ := testQueries.CreateCustomer(context.Background(), util.GenerateRandomString(9))
	return customer
}

func createRandomBatch() Batch {
	bundle := createRandomBundle()
	customer := createRandomCustomer()
	param := CreateBatchParams{
		Name:             "Random batch",
		ActivationStatus: "Active",
		CustomerID:       customer.ID,
		MrcID:            int64(bundle.ID),
	}
	batch, _ := testQueries.CreateBatch(context.Background(), param)
	return batch
}

func TestCreateBatch(t *testing.T) {
	bundle := createRandomBundle()
	customer := createRandomCustomer()

	param := CreateBatchParams{
		Name:             "Random batch",
		ActivationStatus: "Active",
		CustomerID:       customer.ID,
		MrcID:            int64(bundle.ID),
	}

	batch, err := testQueries.CreateBatch(context.Background(), param)

	require.NoError(t, err)
	require.Equal(t, batch.Name, param.Name)
	require.Equal(t, batch.ActivationStatus, param.ActivationStatus)
	require.Equal(t, batch.CustomerID, param.CustomerID)
	require.Equal(t, batch.MrcID, param.MrcID)

}

func TestDeleteBatch(t *testing.T) {
	batch := createRandomBatch()
	err := testQueries.DeleteBatch(context.Background(), batch.ID)
	require.NoError(t, err)

	batch_2, err := testQueries.GetBatchForUpdate(context.Background(), batch.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, batch_2)
}

func TestGetBatch(t *testing.T) {
	batch := createRandomBatch()

	batch_2, err := testQueries.GetBatchForUpdate(context.Background(), batch.ID)

	require.NoError(t, err)
	require.Equal(t, batch, batch_2)

}

func TestUpdateBatch(t *testing.T) {
	batch := createRandomBatch()
	new_bundle := createRandomBundle()
	customer := createRandomCustomer()

	arg := UpdateBatchParams{
		ID:               batch.ID,
		MrcID:            new_bundle.ID,
		CustomerID:       customer.ID,
		ActivationStatus: "Inactive",
		NoOfDevices:      batch.NoOfDevices + 10,
		DeliveryDate:     sql.NullTime{Time: time.Now(), Valid: true},
		WarrantyEnd:      sql.NullTime{Time: time.Now().Add(100), Valid: true},
	}

	batch_2, err := testQueries.UpdateBatch(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, batch_2.ID, arg.ID)
	require.Equal(t, batch_2.MrcID, arg.MrcID)
	require.Equal(t, batch_2.CustomerID, arg.CustomerID)
	require.Equal(t, batch_2.ActivationStatus, arg.ActivationStatus)
	require.Equal(t, batch_2.NoOfDevices, arg.NoOfDevices)
	require.WithinDuration(t, batch_2.DeliveryDate.Time, arg.DeliveryDate.Time, time.Second)
	require.WithinDuration(t, batch_2.WarrantyEnd.Time, arg.WarrantyEnd.Time, time.Second)
}

func TestListAllBatches(t *testing.T) {
	var batch Batch
	for x := 0; x < 5; x++ {
		batch = createRandomBatch()
	}
	arg := ListAllBatchesParams{
		Limit:  5,
		Offset: 0,
		Name:   sql.NullString{String: batch.Name, Valid: true},
	}
	batch_2, err := testQueries.ListAllBatches(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, batch_2, 5)
}
