package orders

import "sushi-mart/internal/database"

// service for handling customer orders
type OrderService interface {
	AddOrder() error
	CancelOrder() error
	GetOrders() error
	GetStatus() error
}

type OrderServiceImpl struct {
	Queries *database.Queries
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

/*
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
*/
