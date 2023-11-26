package db

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCharge() Charge {
	customer := createRandomCustomer()
	arg := CreateChargesParams{
		Paid:       rand.Float64() * 1000,
		Due:        rand.Float64() * 1000,
		CustomerID: customer.ID,
	}
	charge, _ := testQueries.CreateCharges(context.Background(), arg)
	return charge
}
func TestCreateCharges(t *testing.T) {
	customer := createRandomCustomer()
	arg := CreateChargesParams{
		Paid:       rand.Float64() * 1000,
		Due:        rand.Float64() * 1000,
		CustomerID: customer.ID,
	}
	charge, err := testQueries.CreateCharges(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, charge.Paid, arg.Paid)
	require.Equal(t, charge.Due, arg.Due)
	require.Equal(t, charge.CustomerID, arg.CustomerID)
}

func TestDeleteCharges(t *testing.T) {
	charges := createRandomCharge()
	err := testQueries.DeleteCharges(context.Background(), charges.ID)
	require.NoError(t, err)
}

func TestUpdateCharges(t *testing.T) {
	charges := createRandomCharge()
	arg := UpdateChargesParams{
		CustomerID: charges.CustomerID,
		Paid:       charges.Paid + 1.3,
		Due:        charges.Due - 1.4,
	}
	charges_2, err := testQueries.UpdateCharges(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, charges.Paid+1.3, charges_2.Paid)
	require.Equal(t, charges.CustomerID, charges_2.CustomerID)
	require.Equal(t, charges.Due-1.4, charges_2.Due)
}
