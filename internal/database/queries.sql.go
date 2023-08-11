// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: queries.sql

package database

import (
	"context"
	"database/sql"
	"time"
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
	UnitPrice    string
	DateAdded    time.Time
	DateModified time.Time
}

func (q *Queries) AddProduct(ctx context.Context, arg AddProductParams) error {
	_, err := q.db.ExecContext(ctx, addProduct,
		arg.Name,
		arg.Quantity,
		arg.Category,
		arg.UnitPrice,
		arg.DateAdded,
		arg.DateModified,
	)
	return err
}

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers (
  username, password, email, phone, address
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, username, password, email, phone, address
`

type CreateCustomerParams struct {
	Username string
	Password string
	Email    string
	Phone    sql.NullString
	Address  sql.NullString
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.Phone,
		arg.Address,
	)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Phone,
		&i.Address,
	)
	return i, err
}

const createWallet = `-- name: CreateWallet :exec
INSERT INTO wallet(
  balance, wallet_type, date_added, date_modified, customer_id
) VALUES(
  $1, $2, $3, $4, $5
)
`

type CreateWalletParams struct {
	Balance      string
	WalletType   string
	DateAdded    time.Time
	DateModified time.Time
	CustomerID   sql.NullInt32
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) error {
	_, err := q.db.ExecContext(ctx, createWallet,
		arg.Balance,
		arg.WalletType,
		arg.DateAdded,
		arg.DateModified,
		arg.CustomerID,
	)
	return err
}

const deletProduct = `-- name: DeletProduct :exec
DELETE FROM productItems WHERE id=$1
`

func (q *Queries) DeletProduct(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletProduct, id)
	return err
}

const getAllProducts = `-- name: GetAllProducts :many
SELECT id, name, quantity, category, unit_price, date_added, date_modified FROM productItems
`

func (q *Queries) GetAllProducts(ctx context.Context) ([]Productitem, error) {
	rows, err := q.db.QueryContext(ctx, getAllProducts)
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
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCustomer = `-- name: GetCustomer :one
SELECT id, username, password, email, phone, address FROM customers
WHERE email = $1
`

func (q *Queries) GetCustomer(ctx context.Context, email string) (Customer, error) {
	row := q.db.QueryRowContext(ctx, getCustomer, email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Phone,
		&i.Address,
	)
	return i, err
}

const getWallet = `-- name: GetWallet :one
SELECT c.username, w.balance, w.wallet_type FROM wallet w 
INNER JOIN customers c ON w.customer_id=c.id
WHERE c.id = $1
`

type GetWalletRow struct {
	Username   string
	Balance    string
	WalletType string
}

func (q *Queries) GetWallet(ctx context.Context, id int32) (GetWalletRow, error) {
	row := q.db.QueryRowContext(ctx, getWallet, id)
	var i GetWalletRow
	err := row.Scan(&i.Username, &i.Balance, &i.WalletType)
	return i, err
}

const updateBalance = `-- name: UpdateBalance :exec
UPDATE wallet 
SET balance = CASE WHEN $1::boolean THEN $2::DECIMAL(10,3) ELSE balance END,
    wallet_type = CASE WHEN $3::boolean THEN $4::VARCHAR(20) ELSE wallet_type END,
    date_modified = CASE WHEN $5::boolean THEN $6::TIMESTAMP ELSE date_modified END
WHERE id = $7
`

type UpdateBalanceParams struct {
	UpdateBalance      bool
	Balance            string
	UpdateWalletType   bool
	WalletType         string
	UpdateDateModified bool
	DateModified       time.Time
	ID                 int32
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) error {
	_, err := q.db.ExecContext(ctx, updateBalance,
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
    unit_price = CASE WHEN $7::boolean THEN $8::DECIMAL(5,2) ELSE unit_price END,
    date_modified = CASE WHEN $9::boolean THEN $10::TIMESTAMP ELSE date_modified END
WHERE id = $11
RETURNING id, name, quantity, category, unit_price, date_added, date_modified
`

type UpdateProductParams struct {
	UpdateName         bool
	Name               string
	UpdateQuantity     bool
	Quantity           int32
	UpdateCategory     bool
	Category           string
	UpdateUnitPrice    bool
	UnitPrice          string
	UpdateDateModified bool
	DateModified       time.Time
	ID                 int32
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Productitem, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
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
	)
	return i, err
}
