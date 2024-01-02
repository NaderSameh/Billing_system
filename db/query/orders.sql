-- name: CreateOrder :one
INSERT INTO orders (
  start_date, end_date,batch_id,bundle_id,nrc
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;


-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;



-- name: ListAllOrders :many
SELECT * FROM orders;

-- name: ListOrdersByBundleID :many
SELECT * FROM orders
WHERE bundle_id = $1
ORDER BY bundle_id;


-- name: ListOrdersByBatchID :many
SELECT * FROM orders
WHERE batch_id = $1
ORDER BY batch_id;


-- name: ListAllActiveOrders :many
SELECT * FROM orders
WHERE start_date <= NOW()
AND end_date > NOW()
AND nrc is NULL;


-- name: UpdateOrders :one
UPDATE orders
SET
  nrc = COALESCE(sqlc.narg('nrc'), nrc),
  bundle_id = $2,
  start_date = COALESCE(sqlc.narg('start_date'), start_date),
  end_date = $3
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
