-- name: CreateOrder :one
INSERT INTO orders (
  start_date, end_date,batch_id,bundle_id,nrc
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;



-- name: ListOrdersByBundleID :many
SELECT * FROM orders
WHERE bundle_id = $1
ORDER BY bundle_id;


-- name: ListOrdersByBatchID :many
SELECT * FROM orders
WHERE batch_id = $1
ORDER BY batch_id;


-- name: UpdateOrders :one
UPDATE orders
SET nrc = $2,
bundle_id = $3,
start_date = $4,
end_date = $5
WHERE id = $1
RETURNING *;


-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
