-- name: CreateBatch :one
INSERT INTO batches (
  name, activation_status, customer_id, 
  no_of_devices, delivery_date, warranty_end
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;


-- name: GetBatchForUpdate :one
SELECT * FROM batches
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: GetBatchByName :one
SELECT * FROM batches
WHERE name = $1 LIMIT 1;


-- name: ListAllBatches :many
SELECT * FROM batches
WHERE (name = sqlc.narg('name') OR sqlc.narg('name') IS NULL)
AND (customer_id = sqlc.narg('customer_id') OR sqlc.narg('customer_id') IS NULL)
ORDER BY id
LIMIT $1
OFFSET $2;



-- name: UpdateBatch :one
UPDATE batches
SET
customer_id = $3,
activation_status = $4,
no_of_devices = $5,
delivery_date = $6,
warranty_end = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBatch :exec
DELETE FROM batches
WHERE id = $1;
