-- name: CreateCharges :one
INSERT INTO charges (
  paid, due, customer_id
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: UpdateCharges :one
UPDATE charges
SET paid = $2,
due = $3
WHERE customer_id = $1
RETURNING *;

-- name: DeleteCharges :exec
DELETE FROM charges
WHERE customer_id = $1;
