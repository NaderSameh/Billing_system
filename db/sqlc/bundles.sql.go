// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: bundles.sql

package db

import (
	"context"
	"encoding/json"
)

const addBundleToCustomer = `-- name: AddBundleToCustomer :exec
INSERT INTO bundles_customers (
  bundles_id, customers_id
) VALUES (
  $1, $2
)
`

type AddBundleToCustomerParams struct {
	BundlesID   int64 `json:"bundles_id"`
	CustomersID int64 `json:"customers_id"`
}

func (q *Queries) AddBundleToCustomer(ctx context.Context, arg AddBundleToCustomerParams) error {
	_, err := q.db.ExecContext(ctx, addBundleToCustomer, arg.BundlesID, arg.CustomersID)
	return err
}

const createBundle = `-- name: CreateBundle :one
INSERT INTO bundles (
  mrc, description
) VALUES (
  $1, $2
)
RETURNING id, mrc, description
`

type CreateBundleParams struct {
	Mrc         float64 `json:"mrc"`
	Description string  `json:"description"`
}

func (q *Queries) CreateBundle(ctx context.Context, arg CreateBundleParams) (Bundle, error) {
	row := q.db.QueryRowContext(ctx, createBundle, arg.Mrc, arg.Description)
	var i Bundle
	err := row.Scan(&i.ID, &i.Mrc, &i.Description)
	return i, err
}

const deleteBundle = `-- name: DeleteBundle :exec
DELETE FROM bundles
WHERE id = $1
`

func (q *Queries) DeleteBundle(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBundle, id)
	return err
}

const getBundleByID = `-- name: GetBundleByID :one
SELECT id, mrc, description FROM bundles 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBundleByID(ctx context.Context, id int64) (Bundle, error) {
	row := q.db.QueryRowContext(ctx, getBundleByID, id)
	var i Bundle
	err := row.Scan(&i.ID, &i.Mrc, &i.Description)
	return i, err
}

const listAllBundles = `-- name: ListAllBundles :many
SELECT id, mrc, description FROM bundles
ORDER BY id
`

func (q *Queries) ListAllBundles(ctx context.Context) ([]Bundle, error) {
	rows, err := q.db.QueryContext(ctx, listAllBundles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bundle{}
	for rows.Next() {
		var i Bundle
		if err := rows.Scan(&i.ID, &i.Mrc, &i.Description); err != nil {
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

const listBundlesByCustomerID = `-- name: ListBundlesByCustomerID :many
SELECT bundles.id, bundles.mrc, bundles.description
FROM bundles
JOIN bundles_customers ON bundles.id = bundles_customers.bundles_id
WHERE bundles_customers.customers_id = $1
ORDER BY bundles_customers.customers_id
`

func (q *Queries) ListBundlesByCustomerID(ctx context.Context, customersID int64) ([]Bundle, error) {
	rows, err := q.db.QueryContext(ctx, listBundlesByCustomerID, customersID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bundle{}
	for rows.Next() {
		var i Bundle
		if err := rows.Scan(&i.ID, &i.Mrc, &i.Description); err != nil {
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

const listBundlesWithCustomer = `-- name: ListBundlesWithCustomer :many
WITH BundleCustomers AS (
  SELECT bc.bundles_id, json_agg(json_build_object('customer_id', c.id, 'customer', c.customer)) AS assigned_customers
  FROM bundles_customers bc
  JOIN customers c ON c.id = bc.customers_id
  GROUP BY bc.bundles_id
)
SELECT b.id AS bundle_id, b.mrc, b.description, COALESCE(bc.assigned_customers, '[]'::json)
FROM bundles b
LEFT JOIN BundleCustomers bc ON b.id = bc.bundles_id
`

type ListBundlesWithCustomerRow struct {
	BundleID          int64           `json:"bundle_id"`
	Mrc               float64         `json:"mrc"`
	Description       string          `json:"description"`
	AssignedCustomers json.RawMessage `json:"assigned_customers"`
}

func (q *Queries) ListBundlesWithCustomer(ctx context.Context) ([]ListBundlesWithCustomerRow, error) {
	rows, err := q.db.QueryContext(ctx, listBundlesWithCustomer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListBundlesWithCustomerRow{}
	for rows.Next() {
		var i ListBundlesWithCustomerRow
		if err := rows.Scan(
			&i.BundleID,
			&i.Mrc,
			&i.Description,
			&i.AssignedCustomers,
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
