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


-- name: ListBundlesWithCustomer :many
WITH BundleCustomers AS (
  SELECT bc.bundles_id, json_agg(json_build_object('customer_id', c.id, 'customer', c.customer)) AS assigned_customers
  FROM bundles_customers bc
  JOIN customers c ON c.id = bc.customers_id
  GROUP BY bc.bundles_id
)
SELECT b.id AS bundle_id, b.mrc, b.description, COALESCE(bc.assigned_customers, '[]'::json)
FROM bundles b
LEFT JOIN BundleCustomers bc ON b.id = bc.bundles_id;



-- name: DeleteOldBundleCustomers :exec
DELETE FROM bundles_customers
WHERE bundles_id = $1;

-- name: InsertNewBundleCustomers :exec
INSERT INTO bundles_customers (bundles_id, customers_id)
SELECT
    $2 AS bundles_id,
    (data->>'customer_id')::bigint AS customers_id
FROM json_array_elements($1::json) AS data;




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

