// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"

)

type Querier interface {
	AddBundleToCustomer(ctx context.Context, arg AddBundleToCustomerParams) error
	CreateBatch(ctx context.Context, arg CreateBatchParams) (Batch, error)
	CreateBundle(ctx context.Context, arg CreateBundleParams) (Bundle, error)
	CreateCharges(ctx context.Context, arg CreateChargesParams) (Charge, error)
	CreateCustomer(ctx context.Context, customer string) (Customer, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (PaymentLog, error)
	DeleteBatch(ctx context.Context, id int64) error
	DeleteBundle(ctx context.Context, id int64) error
	DeleteCharges(ctx context.Context, customerID int64) error
	DeleteOrder(ctx context.Context, id int64) error
	DeletePayment(ctx context.Context, id int64) error
	GetBatchForUpdate(ctx context.Context, id int64) (Batch, error)
	GetPaymentForUpdate(ctx context.Context, id int64) (PaymentLog, error)
	ListAllBatches(ctx context.Context, arg ListAllBatchesParams) ([]Batch, error)
	ListBundlesByCustomerID(ctx context.Context, customersID int64) ([]ListBundlesByCustomerIDRow, error)
	ListOrdersByBatchID(ctx context.Context, batchID int64) ([]Order, error)
	ListOrdersByBundleID(ctx context.Context, bundleID int64) ([]Order, error)
	ListPaymentByConfirmation(ctx context.Context, confirmed bool) ([]PaymentLog, error)
	UpdateBatch(ctx context.Context, arg UpdateBatchParams) (Batch, error)
	UpdateCharges(ctx context.Context, arg UpdateChargesParams) (Charge, error)
	UpdateOrders(ctx context.Context, arg UpdateOrdersParams) (Order, error)
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (PaymentLog, error)
}

var _ Querier = (*Queries)(nil)
