package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func TestCreateCustomer(t *testing.T) {

	name := util.GenerateRandomString(9)
	customer2, err := testQueries.CreateCustomer(context.Background(), name)
	require.NoError(t, err)
	require.Equal(t, customer2.Paid, float64(0))
	require.Equal(t, customer2.Due, float64(0))
	require.Equal(t, customer2.Customer, name)
	_, err = testQueries.CreateCustomer(context.Background(), name)
	require.Error(t, err)
	require.ErrorContains(t, err, "unique")

}

func TestAddToDue(t *testing.T) {
	customer := createRandomCustomer()
	arg := AddToDueParams{
		Amount: 10,
		ID:     customer.ID,
	}
	customer_2, err := testQueries.AddToDue(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, customer.Due+10, customer_2.Due)

	arg = AddToDueParams{
		Amount: -customer_2.Due,
		ID:     customer.ID,
	}
	customer_2, err = testQueries.AddToDue(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, float64(0), customer_2.Due)
}

func TestListAll(t *testing.T) {

	customers, err := testQueries.ListAllCharges(context.Background(), sql.NullString{String: "randomdada", Valid: true})
	require.NoError(t, err)
	require.Len(t, customers, 0)
}
