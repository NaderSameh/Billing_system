// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: payment_logs.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payment_logs (
  payment, due_date, order_id, confirmed, customer_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, payment, due_date, confirmation_date, order_id, confirmed, customer_id
`

type CreatePaymentParams struct {
	Payment    float64   `json:"payment"`
	DueDate    time.Time `json:"due_date"`
	OrderID    int64     `json:"order_id"`
	Confirmed  bool      `json:"confirmed"`
	CustomerID int64     `json:"customer_id"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (PaymentLog, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.Payment,
		arg.DueDate,
		arg.OrderID,
		arg.Confirmed,
		arg.CustomerID,
	)
	var i PaymentLog
	err := row.Scan(
		&i.ID,
		&i.Payment,
		&i.DueDate,
		&i.ConfirmationDate,
		&i.OrderID,
		&i.Confirmed,
		&i.CustomerID,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payment_logs
WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePayment, id)
	return err
}

const getPaymentForUpdate = `-- name: GetPaymentForUpdate :one
SELECT id, payment, due_date, confirmation_date, order_id, confirmed, customer_id FROM payment_logs
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetPaymentForUpdate(ctx context.Context, id int64) (PaymentLog, error) {
	row := q.db.QueryRowContext(ctx, getPaymentForUpdate, id)
	var i PaymentLog
	err := row.Scan(
		&i.ID,
		&i.Payment,
		&i.DueDate,
		&i.ConfirmationDate,
		&i.OrderID,
		&i.Confirmed,
		&i.CustomerID,
	)
	return i, err
}

const listAllPaymentsCount = `-- name: ListAllPaymentsCount :one
SELECT COUNT(*) FROM payment_logs
WHERE (confirmed = $1 OR $1 IS NULL)
AND (customer_id = $2 OR $2 IS NULL)
`

type ListAllPaymentsCountParams struct {
	Confirmed  sql.NullBool  `json:"confirmed"`
	CustomerID sql.NullInt64 `json:"customer_id"`
}

func (q *Queries) ListAllPaymentsCount(ctx context.Context, arg ListAllPaymentsCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, listAllPaymentsCount, arg.Confirmed, arg.CustomerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listPayments = `-- name: ListPayments :many
SELECT id, payment, due_date, confirmation_date, order_id, confirmed, customer_id FROM payment_logs
WHERE (confirmed = $3 OR $3 IS NULL) AND
(customer_id = $4 OR $4 IS NULL)
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPaymentsParams struct {
	Limit      int32         `json:"limit"`
	Offset     int32         `json:"offset"`
	Confirmed  sql.NullBool  `json:"confirmed"`
	CustomerID sql.NullInt64 `json:"customer_id"`
}

func (q *Queries) ListPayments(ctx context.Context, arg ListPaymentsParams) ([]PaymentLog, error) {
	rows, err := q.db.QueryContext(ctx, listPayments,
		arg.Limit,
		arg.Offset,
		arg.Confirmed,
		arg.CustomerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PaymentLog{}
	for rows.Next() {
		var i PaymentLog
		if err := rows.Scan(
			&i.ID,
			&i.Payment,
			&i.DueDate,
			&i.ConfirmationDate,
			&i.OrderID,
			&i.Confirmed,
			&i.CustomerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePayment = `-- name: UpdatePayment :one
UPDATE payment_logs
SET   due_date = COALESCE($4, due_date),
confirmed = $2,
confirmation_date = $3
WHERE id = $1
RETURNING id, payment, due_date, confirmation_date, order_id, confirmed, customer_id
`

type UpdatePaymentParams struct {
	ID               int64        `json:"id"`
	Confirmed        bool         `json:"confirmed"`
	ConfirmationDate sql.NullTime `json:"confirmation_date"`
	DueDate          sql.NullTime `json:"due_date"`
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (PaymentLog, error) {
	row := q.db.QueryRowContext(ctx, updatePayment,
		arg.ID,
		arg.Confirmed,
		arg.ConfirmationDate,
		arg.DueDate,
	)
	var i PaymentLog
	err := row.Scan(
		&i.ID,
		&i.Payment,
		&i.DueDate,
		&i.ConfirmationDate,
		&i.OrderID,
		&i.Confirmed,
		&i.CustomerID,
	)
	return i, err
}
