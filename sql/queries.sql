-- name: CreateCustomer :exec
INSERT INTO customers (
  username, password, email, phone, address
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: GetCustomer :one
SELECT * FROM customers
WHERE email = $1 and is_active = TRUE;


-- name: GetAllProducts :many
SELECT * FROM productItems;

-- name: GetProductItem :one
SELECT * FROM productItems WHERE id = $1 and is_active = TRUE;

-- name: AddProduct :exec
INSERT INTO productItems(
  name, quantity, category, unit_price, date_added, date_modified
) VALUES(
  $1, $2, $3, $4, $5, $6
);

-- name: DeletProduct :execrows
UPDATE productItems set is_active = $2 
WHERE id = $1;

-- name: UpdateProduct :one
UPDATE productItems
SET name = CASE WHEN @update_name::boolean THEN @name::VARCHAR(50) ELSE name END,
    quantity = CASE WHEN @update_quantity::boolean THEN @quantity::INT ELSE quantity END,
    category = CASE WHEN @update_category::boolean THEN @category::VARCHAR(50) ELSE category END,
    unit_price = CASE WHEN @update_unit_price::boolean THEN @unit_price::DECIMAL(10,2) ELSE unit_price END,
    date_modified = CASE WHEN @update_date_modified::boolean THEN @date_modified::TIMESTAMP ELSE date_modified END
WHERE id = @id and is_active = TRUE
RETURNING *;  

-- name: CreateWallet :exec
INSERT INTO wallet(
  balance, wallet_type, date_added, date_modified, customer_id
) VALUES(
  $1, $2, $3, $4, $5
);

-- name: GetWallet :one
SELECT c.username, w.balance, w.wallet_type, w.date_added FROM wallet w 
INNER JOIN customers c ON w.customer_id=c.id
WHERE c.id = $1 and is_active = TRUE;

-- name: UpdateBalance :exec
UPDATE wallet 
SET balance = CASE WHEN @update_balance::boolean THEN @balance::DECIMAL(20,3) ELSE balance END,
    wallet_type = CASE WHEN @update_wallet_type::boolean THEN @wallet_type::VARCHAR(20) ELSE wallet_type END,
    date_modified = CASE WHEN @update_date_modified::boolean THEN @date_modified::TIMESTAMP ELSE date_modified END
WHERE id = @id and is_active = TRUE;

-- name: AddReview :exec
INSERT INTO productReviews(
  rating, review_text, review_date, customer_id, product_id
) VALUES (
  $1, $2, $3, $4, $5
);


-- name: GetAvgCustomerRatings :many
SELECT p.name, p.category, ROUND(COALESCE(AVG(pr.rating),0)::numeric,2) AS average_rating FROM productItems p 
LEFT JOIN productReviews pr ON p.id = pr.product_id 
GROUP BY p.id;

-- name: GetMostOrdersPlaced :many
SELECT c.username, c.email, COUNT(o.id) as orders_count FROM customers c
INNER JOIN orders o ON c.id = o.customer_id 
GROUP BY c.id 
ORDER BY orders_count DESC;


-- name: PlaceOrder :exec
INSERT INTO orders(
  order_status, total_amt, units, payment_type, order_date, customer_id, product_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: CancelOrder :execrows
UPDATE orders SET order_status = $3, is_active = FALSE
WHERE id = $1 AND order_status = $2 and is_active = TRUE;

-- name: DeliverOrder :execrows
UPDATE orders SET order_status = $2
WHERE order_status = $1 and is_active = TRUE;

-- name: GetAllPlacedOrders :many
SELECT o.id as order_id, o.order_date, o.order_status, o.total_amt, c.username, p.name as product_name
FROM orders o INNER JOIN customers c ON o.customer_id = c.id
INNER JOIN productItems p ON o.product_id = p.id
WHERE c.id = $1 
ORDER BY o.order_date DESC;

-- name: ValidateProductOrderReview :one
SELECT * from orders 
WHERE product_id = $1 and is_active = TRUE;