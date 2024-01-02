-- name: CreatePayment :one
INSERT INTO payment_logs (
  payment, due_date, order_id, confirmed, customer_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;



-- name: ListPayments :many
SELECT * FROM payment_logs
WHERE (confirmed = sqlc.narg('confirmed') OR sqlc.narg('confirmed') IS NULL) AND
(customer_id = sqlc.narg('customer_id') OR sqlc.narg('customer_id') IS NULL)
ORDER BY id
LIMIT $1
OFFSET $2;



-- name: GetPaymentForUpdate :one
SELECT * FROM payment_logs
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: ListAllPaymentsCount :one
SELECT COUNT(*) FROM payment_logs
WHERE (confirmed = sqlc.narg('confirmed') OR sqlc.narg('confirmed') IS NULL)
AND (customer_id = sqlc.narg('customer_id') OR sqlc.narg('customer_id') IS NULL);



-- name: UpdatePayment :one
UPDATE payment_logs
SET   due_date = COALESCE(sqlc.narg('due_date'), due_date),
confirmed = $2,
confirmation_date = $3
WHERE id = $1
RETURNING *;


-- name: DeletePayment :exec
DELETE FROM payment_logs
WHERE id = $1;
