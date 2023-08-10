-- name: CreateCustomer :one
INSERT INTO customers (
  username, password, email, phone, address
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE email = $1;
