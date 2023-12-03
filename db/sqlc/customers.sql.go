// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: customers.sql

package db

import (
	"context"
	"database/sql"
)

const addToDue = `-- name: AddToDue :one
UPDATE customers
SET due = due + $1
WHERE id = $2
RETURNING id, customer, paid, due
`

type AddToDueParams struct {
	Amount float64 `json:"amount"`
	ID     int64   `json:"id"`
}

func (q *Queries) AddToDue(ctx context.Context, arg AddToDueParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, addToDue, arg.Amount, arg.ID)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Customer,
		&i.Paid,
		&i.Due,
	)
	return i, err
}

const addToPaid = `-- name: AddToPaid :one
UPDATE customers
SET paid = paid + $1,
due = due - $1
WHERE id = $2
RETURNING id, customer, paid, due
`

type AddToPaidParams struct {
	Amount float64 `json:"amount"`
	ID     int64   `json:"id"`
}

func (q *Queries) AddToPaid(ctx context.Context, arg AddToPaidParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, addToPaid, arg.Amount, arg.ID)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Customer,
		&i.Paid,
		&i.Due,
	)
	return i, err
}

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers (
  customer, due, paid
) VALUES (
  $1, 0, 0 
)
RETURNING id, customer, paid, due
`

func (q *Queries) CreateCustomer(ctx context.Context, customer string) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer, customer)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Customer,
		&i.Paid,
		&i.Due,
	)
	return i, err
}

const getCustomerID = `-- name: GetCustomerID :one
SELECT id, customer, paid, due FROM customers
WHERE customer = $1 LIMIT 1
`

func (q *Queries) GetCustomerID(ctx context.Context, customer string) (Customer, error) {
	row := q.db.QueryRowContext(ctx, getCustomerID, customer)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Customer,
		&i.Paid,
		&i.Due,
	)
	return i, err
}

const listAllCharges = `-- name: ListAllCharges :many
SELECT id, customer, paid, due FROM customers
WHERE (customer = $1 OR $1 IS NULL)
ORDER BY id
`

func (q *Queries) ListAllCharges(ctx context.Context, name sql.NullString) ([]Customer, error) {
	rows, err := q.db.QueryContext(ctx, listAllCharges, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Customer{}
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.Customer,
			&i.Paid,
			&i.Due,
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
