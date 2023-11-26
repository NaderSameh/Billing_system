package db

import (
	"context"
	"math/rand"
	"testing"

	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func TestCreateBundle(t *testing.T) {
	arg := CreateBundleParams{
		Mrc:         rand.Int31n(100),
		Description: util.GenerateRandomString(25),
	}
	bundle, err := testQueries.CreateBundle(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, bundle.Mrc, arg.Mrc)
	require.Equal(t, bundle.Description, arg.Description)
}

func TestDeleteBundle(t *testing.T) {
	batch := createRandomBundle()
	err := testQueries.DeleteBundle(context.Background(), batch.ID)
	require.NoError(t, err)
}

func TestAddBundleToCustomer(t *testing.T) {
	bundle := createRandomBundle()
	customer := createRandomCustomer()
	arg := AddBundleToCustomerParams{
		BundlesID:   bundle.ID,
		CustomersID: customer.ID,
	}

	err := testQueries.AddBundleToCustomer(context.Background(), arg)
	require.NoError(t, err)
}

func TestListAllBundles(t *testing.T) {
	Customer := createRandomCustomer()
	var bundle Bundle
	for n := 0; n < 10; n++ {
		bundle = createRandomBundle()
		arg := AddBundleToCustomerParams{
			BundlesID:   bundle.ID,
			CustomersID: Customer.ID,
		}
		testQueries.AddBundleToCustomer(context.Background(), arg)
	}
	bundles, err := testQueries.ListBundlesByCustomerID(context.Background(), Customer.ID)
	require.NoError(t, err)
	require.Len(t, bundles, 10)
}
