package inventory

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"
)

// service for product items
type InventoryService interface {
	GetAllProducts(context.Context) (*GetAllProductsResp, *common.ErrorResponse)
	AddProduct(context.Context, *AddProductReq) *common.ErrorResponse
	DeleteProduct(context.Context, int) *common.ErrorResponse
	UpdateProduct(context.Context, int, *UpdateProductReq) (*ProductResp, *common.ErrorResponse)
}

type InventoryServiceImpl struct {
	Queries database.Querier
}

type Cache struct {
	InventoryService
}
type Validator struct {
	InventoryService
}

type RoutesWrapper struct {
	InventoryService
}

func New(Queries *database.Queries) *RoutesWrapper {
	return &RoutesWrapper{
		InventoryService: &Validator{
			InventoryService: &Cache{
				InventoryService: &InventoryServiceImpl{
					Queries: Queries,
				},
			},
		},
	}
}
