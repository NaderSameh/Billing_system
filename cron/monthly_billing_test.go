package cron

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/billing_system/db/mock"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestUpdateMonthlyBillingMockSecond(t *testing.T) {
	n_ok := 10
	n_empty := 0
	testCases := []struct {
		name       string
		n          int
		buildstuds func(store *mockdb.MockStore)
	}{
		{
			name: "OK",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Bundle{}, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Customer{}, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.PaymentLog{}, nil)
			},
		},

		{
			name: "No orders",
			n:    n_empty,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_empty)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(n_empty).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(n_empty).Return(db.Bundle{}, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(n_empty).Return(db.Customer{}, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(n_empty).Return(db.PaymentLog{}, nil)
			},
		},

		{
			name: "interal server error - list active",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, sql.ErrConnDone)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(0).Return(db.Batch{}, nil)
			},
		},
		{
			name: "interal server error - get batch",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(db.Batch{}, sql.ErrConnDone)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(0).Return(db.Bundle{}, nil)
			},
		},
		{
			name: "interal server error - get bundle",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(1).Return(db.Bundle{}, sql.ErrConnDone)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(0).Return(db.Customer{}, nil)
			},
		},
		{
			name: "interal server error - add due",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(1).Return(db.Bundle{}, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(db.Customer{}, sql.ErrConnDone)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(0).Return(db.PaymentLog{}, nil)
			},
		},
		{
			name: "interal server error - create payment",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(1).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(1).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(1).Return(db.Bundle{}, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(1).Return(db.Customer{}, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(1).Return(db.PaymentLog{}, sql.ErrConnDone)
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
			config := CronJobConfiguration{Interval: 1, Unit: CronSecond, Name: "test 1 sec"}
			cron, _ := NewTestCron(t, store, config)
			time.Sleep(time.Millisecond * 500)
			cron.StopCron()

		})
	}
}

func TestUpdateMonthlyBilling(t *testing.T) {
	n_ok := 0
	testCases := []struct {
		name          string
		n             int
		Schedule_date time.Time
		buildstuds    func(store *mockdb.MockStore)
		checkdate     func(schedule_time time.Time)
	}{
		{
			name: "OK",
			n:    n_ok,
			buildstuds: func(store *mockdb.MockStore) {
				orders := make([]db.Order, n_ok)
				store.EXPECT().ListAllActiveOrders(gomock.Any()).Times(n_ok).Return(orders, nil)
				store.EXPECT().GetBatchForUpdate(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Batch{}, nil)
				store.EXPECT().GetBundleByID(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Bundle{}, nil)
				store.EXPECT().AddToDue(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.Customer{}, nil)
				store.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Times(n_ok).Return(db.PaymentLog{}, nil)
			},
			checkdate: func(schedule_time time.Time) {
				now := time.Now()
				estimate_time := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
				estimate_time = estimate_time.AddDate(0, 1, -1)
				require.WithinDuration(t, schedule_time, estimate_time, time.Second)
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
			config := CronJobConfiguration{Interval: 1, Unit: CronEndOfMonth, Name: "Full month test"}
			cron, J := NewTestCron(t, store, config)

			tc.Schedule_date = J.job.ScheduledTime()

			tc.checkdate(tc.Schedule_date)

			cron.StopCron()

		})
	}
}
