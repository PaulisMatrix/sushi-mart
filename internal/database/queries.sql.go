// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: queries.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

const addProduct = `-- name: AddProduct :exec
INSERT INTO productItems(
  name, quantity, category, unit_price, date_added, date_modified
) VALUES(
  $1, $2, $3, $4, $5, $6
)
`

type AddProductParams struct {
	Name         string
	Quantity     int32
	Category     string
	UnitPrice    decimal.Decimal
	DateAdded    pgtype.Timestamp
	DateModified pgtype.Timestamp
}

func (q *Queries) AddProduct(ctx context.Context, arg AddProductParams) error {
	_, err := q.db.Exec(ctx, addProduct,
		arg.Name,
		arg.Quantity,
		arg.Category,
		arg.UnitPrice,
		arg.DateAdded,
		arg.DateModified,
	)
	return err
}

const addReview = `-- name: AddReview :exec
INSERT INTO productReviews(
  rating, review_text, review_date, customer_id, product_id
) VALUES (
  $1, $2, $3, $4, $5
)
`

type AddReviewParams struct {
	Rating     int32
	ReviewText string
	ReviewDate pgtype.Timestamp
	CustomerID pgtype.Int4
	ProductID  pgtype.Int4
}

func (q *Queries) AddReview(ctx context.Context, arg AddReviewParams) error {
	_, err := q.db.Exec(ctx, addReview,
		arg.Rating,
		arg.ReviewText,
		arg.ReviewDate,
		arg.CustomerID,
		arg.ProductID,
	)
	return err
}

const cancelOrder = `-- name: CancelOrder :execrows
UPDATE orders SET order_status = $3, is_active = FALSE
WHERE id = $1 AND order_status = $2 and is_active = TRUE
`

type CancelOrderParams struct {
	ID            int32
	OrderStatus   string
	OrderStatus_2 string
}

func (q *Queries) CancelOrder(ctx context.Context, arg CancelOrderParams) (int64, error) {
	result, err := q.db.Exec(ctx, cancelOrder, arg.ID, arg.OrderStatus, arg.OrderStatus_2)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createCustomer = `-- name: CreateCustomer :exec
INSERT INTO customers (
  username, password, email, phone, address
) VALUES (
  $1, $2, $3, $4, $5
)
`

type CreateCustomerParams struct {
	Username string
	Password string
	Email    string
	Phone    pgtype.Text
	Address  pgtype.Text
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) error {
	_, err := q.db.Exec(ctx, createCustomer,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.Phone,
		arg.Address,
	)
	return err
}

const createWallet = `-- name: CreateWallet :exec
INSERT INTO wallet(
  balance, wallet_type, date_added, date_modified, customer_id
) VALUES(
  $1, $2, $3, $4, $5
)
`

type CreateWalletParams struct {
	Balance      decimal.Decimal
	WalletType   string
	DateAdded    pgtype.Timestamp
	DateModified pgtype.Timestamp
	CustomerID   pgtype.Int4
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) error {
	_, err := q.db.Exec(ctx, createWallet,
		arg.Balance,
		arg.WalletType,
		arg.DateAdded,
		arg.DateModified,
		arg.CustomerID,
	)
	return err
}

const deletProduct = `-- name: DeletProduct :execrows
UPDATE productItems set is_active = $2 
WHERE id = $1
`

type DeletProductParams struct {
	ID       int32
	IsActive bool
}

func (q *Queries) DeletProduct(ctx context.Context, arg DeletProductParams) (int64, error) {
	result, err := q.db.Exec(ctx, deletProduct, arg.ID, arg.IsActive)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deliverOrder = `-- name: DeliverOrder :execrows
UPDATE orders SET order_status = $2
WHERE order_status = $1 and is_active = TRUE
`

type DeliverOrderParams struct {
	OrderStatus   string
	OrderStatus_2 string
}

func (q *Queries) DeliverOrder(ctx context.Context, arg DeliverOrderParams) (int64, error) {
	result, err := q.db.Exec(ctx, deliverOrder, arg.OrderStatus, arg.OrderStatus_2)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllPlacedOrders = `-- name: GetAllPlacedOrders :many
SELECT o.id as order_id, o.order_date, o.order_status, o.total_amt, c.username, p.name as product_name
FROM orders o INNER JOIN customers c ON o.customer_id = c.id
INNER JOIN productItems p ON o.product_id = p.id
WHERE c.id = $1 
ORDER BY o.order_date DESC
`

type GetAllPlacedOrdersRow struct {
	OrderID     int32
	OrderDate   pgtype.Timestamp
	OrderStatus string
	TotalAmt    decimal.Decimal
	Username    string
	ProductName string
}

func (q *Queries) GetAllPlacedOrders(ctx context.Context, id int32) ([]GetAllPlacedOrdersRow, error) {
	rows, err := q.db.Query(ctx, getAllPlacedOrders, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllPlacedOrdersRow
	for rows.Next() {
		var i GetAllPlacedOrdersRow
		if err := rows.Scan(
			&i.OrderID,
			&i.OrderDate,
			&i.OrderStatus,
			&i.TotalAmt,
			&i.Username,
			&i.ProductName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllProducts = `-- name: GetAllProducts :many
SELECT id, name, quantity, category, unit_price, date_added, date_modified, is_active FROM productItems
`

func (q *Queries) GetAllProducts(ctx context.Context) ([]Productitem, error) {
	rows, err := q.db.Query(ctx, getAllProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Productitem
	for rows.Next() {
		var i Productitem
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Quantity,
			&i.Category,
			&i.UnitPrice,
			&i.DateAdded,
			&i.DateModified,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAvgCustomerRatings = `-- name: GetAvgCustomerRatings :many
SELECT p.name, p.category, ROUND(COALESCE(AVG(pr.rating),0)::numeric,2) AS average_rating FROM productItems p 
LEFT JOIN productReviews pr ON p.id = pr.product_id 
GROUP BY p.id
`

type GetAvgCustomerRatingsRow struct {
	Name          string
	Category      string
	AverageRating pgtype.Numeric
}

func (q *Queries) GetAvgCustomerRatings(ctx context.Context) ([]GetAvgCustomerRatingsRow, error) {
	rows, err := q.db.Query(ctx, getAvgCustomerRatings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAvgCustomerRatingsRow
	for rows.Next() {
		var i GetAvgCustomerRatingsRow
		if err := rows.Scan(&i.Name, &i.Category, &i.AverageRating); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCustomer = `-- name: GetCustomer :one
SELECT id, username, password, email, phone, address, is_active FROM customers
WHERE email = $1 and is_active = TRUE
`

func (q *Queries) GetCustomer(ctx context.Context, email string) (Customer, error) {
	row := q.db.QueryRow(ctx, getCustomer, email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Phone,
		&i.Address,
		&i.IsActive,
	)
	return i, err
}

const getMostOrdersPlaced = `-- name: GetMostOrdersPlaced :many
SELECT c.username, c.email, COUNT(o.id) as orders_count FROM customers c
INNER JOIN orders o ON c.id = o.customer_id 
GROUP BY c.id 
ORDER BY orders_count DESC
`

type GetMostOrdersPlacedRow struct {
	Username    string
	Email       string
	OrdersCount int64
}

func (q *Queries) GetMostOrdersPlaced(ctx context.Context) ([]GetMostOrdersPlacedRow, error) {
	rows, err := q.db.Query(ctx, getMostOrdersPlaced)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostOrdersPlacedRow
	for rows.Next() {
		var i GetMostOrdersPlacedRow
		if err := rows.Scan(&i.Username, &i.Email, &i.OrdersCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductItem = `-- name: GetProductItem :one
SELECT id, name, quantity, category, unit_price, date_added, date_modified, is_active FROM productItems WHERE id = $1 and is_active = TRUE
`

func (q *Queries) GetProductItem(ctx context.Context, id int32) (Productitem, error) {
	row := q.db.QueryRow(ctx, getProductItem, id)
	var i Productitem
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Quantity,
		&i.Category,
		&i.UnitPrice,
		&i.DateAdded,
		&i.DateModified,
		&i.IsActive,
	)
	return i, err
}

const getWallet = `-- name: GetWallet :one
SELECT c.username, w.balance, w.wallet_type, w.date_added FROM wallet w 
INNER JOIN customers c ON w.customer_id=c.id
WHERE c.id = $1 and c.is_active = TRUE
`

type GetWalletRow struct {
	Username   string
	Balance    decimal.Decimal
	WalletType string
	DateAdded  pgtype.Timestamp
}

func (q *Queries) GetWallet(ctx context.Context, id int32) (GetWalletRow, error) {
	row := q.db.QueryRow(ctx, getWallet, id)
	var i GetWalletRow
	err := row.Scan(
		&i.Username,
		&i.Balance,
		&i.WalletType,
		&i.DateAdded,
	)
	return i, err
}

const placeOrder = `-- name: PlaceOrder :exec
INSERT INTO orders(
  order_status, total_amt, units, payment_type, order_date, customer_id, product_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
`

type PlaceOrderParams struct {
	OrderStatus string
	TotalAmt    decimal.Decimal
	Units       int32
	PaymentType string
	OrderDate   pgtype.Timestamp
	CustomerID  pgtype.Int4
	ProductID   pgtype.Int4
}

func (q *Queries) PlaceOrder(ctx context.Context, arg PlaceOrderParams) error {
	_, err := q.db.Exec(ctx, placeOrder,
		arg.OrderStatus,
		arg.TotalAmt,
		arg.Units,
		arg.PaymentType,
		arg.OrderDate,
		arg.CustomerID,
		arg.ProductID,
	)
	return err
}

const updateBalance = `-- name: UpdateBalance :exec
UPDATE wallet 
SET balance = CASE WHEN $1::boolean THEN $2::DECIMAL(20,3) ELSE balance END,
    wallet_type = CASE WHEN $3::boolean THEN $4::VARCHAR(20) ELSE wallet_type END,
    date_modified = CASE WHEN $5::boolean THEN $6::TIMESTAMP ELSE date_modified END
WHERE id = $7 and is_active = TRUE
`

type UpdateBalanceParams struct {
	UpdateBalance      bool
	Balance            decimal.Decimal
	UpdateWalletType   bool
	WalletType         string
	UpdateDateModified bool
	DateModified       pgtype.Timestamp
	ID                 int32
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) error {
	_, err := q.db.Exec(ctx, updateBalance,
		arg.UpdateBalance,
		arg.Balance,
		arg.UpdateWalletType,
		arg.WalletType,
		arg.UpdateDateModified,
		arg.DateModified,
		arg.ID,
	)
	return err
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE productItems
SET name = CASE WHEN $1::boolean THEN $2::VARCHAR(50) ELSE name END,
    quantity = CASE WHEN $3::boolean THEN $4::INT ELSE quantity END,
    category = CASE WHEN $5::boolean THEN $6::VARCHAR(50) ELSE category END,
    unit_price = CASE WHEN $7::boolean THEN $8::DECIMAL(10,2) ELSE unit_price END,
    date_modified = CASE WHEN $9::boolean THEN $10::TIMESTAMP ELSE date_modified END
WHERE id = $11 and is_active = TRUE
RETURNING id, name, quantity, category, unit_price, date_added, date_modified, is_active
`

type UpdateProductParams struct {
	UpdateName         bool
	Name               string
	UpdateQuantity     bool
	Quantity           int32
	UpdateCategory     bool
	Category           string
	UpdateUnitPrice    bool
	UnitPrice          decimal.Decimal
	UpdateDateModified bool
	DateModified       pgtype.Timestamp
	ID                 int32
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Productitem, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.UpdateName,
		arg.Name,
		arg.UpdateQuantity,
		arg.Quantity,
		arg.UpdateCategory,
		arg.Category,
		arg.UpdateUnitPrice,
		arg.UnitPrice,
		arg.UpdateDateModified,
		arg.DateModified,
		arg.ID,
	)
	var i Productitem
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Quantity,
		&i.Category,
		&i.UnitPrice,
		&i.DateAdded,
		&i.DateModified,
		&i.IsActive,
	)
	return i, err
}

const validateProductOrderReview = `-- name: ValidateProductOrderReview :one
SELECT id, order_status, total_amt, units, payment_type, order_date, customer_id, product_id, is_active from orders 
WHERE product_id = $1 and is_active = TRUE
`

func (q *Queries) ValidateProductOrderReview(ctx context.Context, productID pgtype.Int4) (Order, error) {
	row := q.db.QueryRow(ctx, validateProductOrderReview, productID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.OrderStatus,
		&i.TotalAmt,
		&i.Units,
		&i.PaymentType,
		&i.OrderDate,
		&i.CustomerID,
		&i.ProductID,
		&i.IsActive,
	)
	return i, err
}
