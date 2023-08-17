package inventory

import (
	"context"
	"net/http"
	"sushi-mart/common"

	"github.com/jackc/pgx/v5"
)

func (v *Validator) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	return v.InventoryService.GetAllProducts(ctx)
}

func (i *InventoryServiceImpl) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithField("method", "GetAllProducts")

	resp, err := i.Queries.GetAllProducts(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithError(err).Error("no products in the inventory yet")
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "no products in the inventory yet",
			}
		}

		logger.WithError(err).Error("error in getting all products from the inventory")
		return nil, &common.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}
	}
	var products []ProductResp

	for _, p := range resp {
		products = append(products, ProductResp{
			Name:         p.Name,
			Quantity:     p.Quantity,
			Category:     p.Category,
			UnitPrice:    p.UnitPrice.Abs().InexactFloat64(),
			DateAdded:    p.DateAdded.Time.Local().String(),
			DateModified: p.DateModified.Time.Local().String(),
		})
	}

	return &GetAllProductsResp{
		Products: products,
	}, nil

}
