-- name: CreateCustomer :one
INSERT INTO customers (
  customer, due, paid
) VALUES (
  $1, 0, 0 
)
RETURNING *;


-- name: GetCustomerID :one
SELECT * FROM customers
WHERE customer = $1 LIMIT 1;



-- name: GetAllCustomers :many
SELECT * FROM customers;

-- name: AddToPaid :one
UPDATE customers
SET paid = paid + sqlc.arg(amount),
due = due - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: AddToDue :one
UPDATE customers
SET due = due + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: ListAllCharges :many
SELECT * FROM customers
WHERE (customer = sqlc.narg('name') OR sqlc.narg('name') IS NULL)
ORDER BY id;

