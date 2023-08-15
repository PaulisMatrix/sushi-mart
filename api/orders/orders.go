package orders

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"
)

// service for handling customer orders
type OrderService interface {
	PlaceOrder(context.Context, *PlaceOrderReq, int) *common.ErrorResponse
	CancelOrder(context.Context, *UpdateOrderReq) *common.ErrorResponse
	GetOrders(context.Context, int) (*GetAllOrdersResp, *common.ErrorResponse)
}

type OrderServiceImpl struct {
	Queries database.Querier
}

type Cache struct {
	OrderService
}
type Validator struct {
	OrderService
}

type RoutesWrapper struct {
	OrderService
}

func New(Queries *database.Queries) *RoutesWrapper {
	return &RoutesWrapper{
		OrderService: &Validator{
			OrderService: &Cache{
				OrderService: &OrderServiceImpl{
					Queries: Queries,
				},
			},
		},
	}
}
