// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"context"
	"database/sql"
)

//go:generate mockgen -destination=./mocks/queries.go -package=mock_queries . Querier
type Querier interface {
	AddProduct(ctx context.Context, arg AddProductParams) error
	AddReview(ctx context.Context, arg AddReviewParams) error
	CancelOrder(ctx context.Context, arg CancelOrderParams) (int64, error)
	CreateCustomer(ctx context.Context, arg CreateCustomerParams) error
	CreateWallet(ctx context.Context, arg CreateWalletParams) error
	DeletProduct(ctx context.Context, arg DeletProductParams) (int64, error)
	DeliverOrder(ctx context.Context, arg DeliverOrderParams) (int64, error)
	GetAllPlacedOrders(ctx context.Context, id int32) ([]GetAllPlacedOrdersRow, error)
	GetAllProducts(ctx context.Context) ([]Productitem, error)
	GetAvgCustomerRatings(ctx context.Context) ([]GetAvgCustomerRatingsRow, error)
	GetCustomer(ctx context.Context, email string) (Customer, error)
	GetMostOrdersPlaced(ctx context.Context) ([]GetMostOrdersPlacedRow, error)
	GetProductItem(ctx context.Context, id int32) (Productitem, error)
	GetWallet(ctx context.Context, id int32) (GetWalletRow, error)
	PlaceOrder(ctx context.Context, arg PlaceOrderParams) error
	UpdateBalance(ctx context.Context, arg UpdateBalanceParams) error
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Productitem, error)
	ValidateProductOrderReview(ctx context.Context, productID sql.NullInt32) (Order, error)
}

var _ Querier = (*Queries)(nil)
