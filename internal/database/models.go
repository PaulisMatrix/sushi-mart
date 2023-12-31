// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type Customer struct {
	ID       int32
	Username string
	Password string
	Email    string
	Phone    pgtype.Text
	Address  pgtype.Text
	IsActive bool
}

type Order struct {
	ID          int32
	OrderStatus string
	TotalAmt    decimal.Decimal
	Units       int32
	PaymentType string
	OrderDate   pgtype.Timestamp
	CustomerID  pgtype.Int4
	ProductID   pgtype.Int4
	IsActive    bool
}

type Productitem struct {
	ID           int32
	Name         string
	Quantity     int32
	Category     string
	UnitPrice    decimal.Decimal
	DateAdded    pgtype.Timestamp
	DateModified pgtype.Timestamp
	IsActive     bool
}

type Productreview struct {
	ID         int32
	Rating     int32
	ReviewText string
	ReviewDate pgtype.Timestamp
	CustomerID pgtype.Int4
	ProductID  pgtype.Int4
	IsActive   bool
}

type Wallet struct {
	ID           int32
	Balance      decimal.Decimal
	WalletType   string
	DateAdded    pgtype.Timestamp
	DateModified pgtype.Timestamp
	CustomerID   pgtype.Int4
	IsActive     bool
}
