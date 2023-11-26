-- name: CreatePayment :one
INSERT INTO payment_logs (
  payment, due_date, order_id, confirmed
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;



-- name: ListPaymentByConfirmation :many
SELECT * FROM payment_logs
WHERE confirmed = $1
ORDER BY id;

-- name: GetPaymentForUpdate :one
SELECT * FROM payment_logs
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;



-- name: UpdatePayment :one
UPDATE payment_logs
SET due_date = $2,
confirmation_date = $3,
order_id = $4,
confirmed = $5
WHERE id = $1
RETURNING *;


-- name: DeletePayment :exec
DELETE FROM payment_logs
WHERE id = $1;
