-- name: CreateCustomer :one
INSERT INTO customers (
  customer
) VALUES (
  $1
)
RETURNING *;


