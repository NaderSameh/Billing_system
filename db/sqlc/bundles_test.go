package db

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func TestCreateBundle(t *testing.T) {
	arg := CreateBundleParams{
		Mrc:         rand.ExpFloat64(),
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

func TestListBundlesWithCustomers(t *testing.T) {
	type bundleWithCustomer struct {
		CustomerID   int64  `json:"customer_id"`
		CustomerName string `json:"customer"`
	}
	var jsonbundleWithCustomer []bundleWithCustomer
	bundle := createRandomBundle()
	Customer := createRandomCustomer()
	Customer2 := createRandomCustomer()
	Customer3 := createRandomCustomer()
	testQueries.AddBundleToCustomer(context.Background(), AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: Customer.ID})
	testQueries.AddBundleToCustomer(context.Background(), AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: Customer3.ID})
	testQueries.AddBundleToCustomer(context.Background(), AddBundleToCustomerParams{BundlesID: bundle.ID, CustomersID: Customer2.ID})
	bundles, err := testQueries.ListBundlesWithCustomer(context.Background())
	require.NoError(t, err)

	for _, b := range bundles {
		if b.BundleID == bundle.ID {
			err := json.Unmarshal(b.AssignedCustomers, &jsonbundleWithCustomer)
			require.NoError(t, err)
			require.Equal(t, b.BundleID, bundle.ID)
			require.Equal(t, jsonbundleWithCustomer[0].CustomerID, Customer.ID)
			require.Equal(t, jsonbundleWithCustomer[0].CustomerName, Customer.Customer)
			require.Equal(t, jsonbundleWithCustomer[1].CustomerID, Customer3.ID)
			require.Equal(t, jsonbundleWithCustomer[1].CustomerName, Customer3.Customer)
			require.Equal(t, jsonbundleWithCustomer[2].CustomerID, Customer2.ID)
			require.Equal(t, jsonbundleWithCustomer[2].CustomerName, Customer2.Customer)

		}
	}

}
