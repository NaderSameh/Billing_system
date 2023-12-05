package cron

import (
	"context"
	"time"

	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/rs/zerolog/log"
)

func (cron *GoCronScheduler) UpdateMonthlyBilling(store db.Store) error {
	active_orders, err := store.ListAllActiveOrders(context.Background())
	if err != nil {
		return err
	}

	for i := range active_orders {

		batch, err := store.GetBatchForUpdate(context.Background(), active_orders[i].BatchID)
		if err != nil {
			return err
		}
		bundle, err := store.GetBundleByID(context.Background(), active_orders[i].BundleID)
		if err != nil {
			return err
		}
		arg := db.AddToDueParams{
			Amount: bundle.Mrc * float64(batch.NoOfDevices),
			ID:     batch.CustomerID,
		}
		_, err = store.AddToDue(context.Background(), arg)
		if err != nil {
			return err
		}
		now := time.Now()
		endOfMonth := time.Date(now.Year(), now.Month()+1, 28, 23, 59, 0, 0, time.UTC)

		arg2 := db.CreatePaymentParams{
			Payment:    arg.Amount,
			OrderID:    active_orders[i].ID,
			CustomerID: batch.CustomerID,
			DueDate:    endOfMonth,
			Confirmed:  false,
		}
		_, err = store.CreatePayment(context.Background(), arg2)
		if err != nil {
			return err
		}
		log.Info().Int64("Order_ID", active_orders[i].ID).Float64("Payment", arg.Amount).Int64("BatchID", batch.ID).Int64("CustomerID", batch.CustomerID).Msg("Cron MRC scheduled task details")
	}

	return nil
}
