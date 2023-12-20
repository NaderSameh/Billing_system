// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddBundleToCustomer(ctx context.Context, arg AddBundleToCustomerParams) error
	AddToDue(ctx context.Context, arg AddToDueParams) (Customer, error)
	AddToPaid(ctx context.Context, arg AddToPaidParams) (Customer, error)
	CreateBatch(ctx context.Context, arg CreateBatchParams) (Batch, error)
	CreateBundle(ctx context.Context, arg CreateBundleParams) (Bundle, error)
	CreateCustomer(ctx context.Context, customer string) (Customer, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (PaymentLog, error)
	DeleteBatch(ctx context.Context, id int64) error
	DeleteBundle(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, id int64) error
	DeletePayment(ctx context.Context, id int64) error
	GetAllCustomers(ctx context.Context) ([]Customer, error)
	GetBatchByName(ctx context.Context, name string) (Batch, error)
	GetBatchForUpdate(ctx context.Context, id int64) (Batch, error)
	GetBundleByID(ctx context.Context, id int64) (Bundle, error)
	GetCustomerID(ctx context.Context, customer string) (Customer, error)
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	GetPaymentForUpdate(ctx context.Context, id int64) (PaymentLog, error)
	ListAllActiveOrders(ctx context.Context) ([]Order, error)
	ListAllBatches(ctx context.Context, arg ListAllBatchesParams) ([]Batch, error)
	ListAllBundles(ctx context.Context) ([]Bundle, error)
	ListAllCharges(ctx context.Context, name sql.NullString) ([]Customer, error)
	ListBundlesByCustomerID(ctx context.Context, customersID int64) ([]Bundle, error)
	ListBundlesWithCustomer(ctx context.Context) ([]ListBundlesWithCustomerRow, error)
	ListOrdersByBatchID(ctx context.Context, batchID int64) ([]Order, error)
	ListOrdersByBundleID(ctx context.Context, bundleID int64) ([]Order, error)
	ListPayments(ctx context.Context, arg ListPaymentsParams) ([]PaymentLog, error)
	UpdateBatch(ctx context.Context, arg UpdateBatchParams) (Batch, error)
	UpdateOrders(ctx context.Context, arg UpdateOrdersParams) (Order, error)
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (PaymentLog, error)
}

var _ Querier = (*Queries)(nil)
