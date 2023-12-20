-- name: CreateBundle :one
INSERT INTO bundles (
  mrc, description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListBundlesByCustomerID :many
SELECT bundles.id, bundles.mrc, bundles.description
FROM bundles
JOIN bundles_customers ON bundles.id = bundles_customers.bundles_id
WHERE bundles_customers.customers_id = $1
ORDER BY bundles_customers.customers_id;


-- name: ListAllBundles :many
SELECT * FROM bundles
ORDER BY id;


-- name: ListBundlesWithCustomers :many
SELECT b.id AS bundle_id, b.mrc, b.description, c.id AS customer_id, c.customer
FROM bundles b
JOIN bundles_customers bc ON b.id = bc.bundles_id
JOIN customers c ON c.id = bc.customers_id
ORDER BY b.id, c.id;


-- name: GetBundleByID :one
SELECT * FROM bundles 
WHERE id = $1 LIMIT 1;


-- name: AddBundleToCustomer :exec
INSERT INTO bundles_customers (
  bundles_id, customers_id
) VALUES (
  $1, $2
);

-- name: DeleteBundle :exec
DELETE FROM bundles
WHERE id = $1;

