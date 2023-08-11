package inventory

import (
	"context"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"
)

func (v *Validator) AddProduct(ctx context.Context, req *AddProductReq) *common.ErrorResponse {
	return v.InventoryService.AddProduct(ctx, req)
}

func (i *InventoryServiceImpl) AddProduct(ctx context.Context, req *AddProductReq) *common.ErrorResponse {

	dbParams := database.AddProductParams{
		Name:         req.Name,
		Quantity:     req.Quantity,
		Category:     req.Category,
		UnitPrice:    strconv.FormatFloat(req.UnitPrice, 'E', -1, 64),
		DateAdded:    time.Now().Local(),
		DateModified: time.Now().Local(),
	}
	err := i.Queries.AddProduct(ctx, dbParams)

	if err != nil {
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
