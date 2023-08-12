package inventory

import (
	"context"
	"net/http"
	"strconv"
	"sushi-mart/common"
)

func (v *Validator) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	return v.InventoryService.GetAllProducts(ctx)
}

func (i *InventoryServiceImpl) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithField("method", "GetAllProducts")

	resp, err := i.Queries.GetAllProducts(ctx)
	if err != nil {
		logger.WithError(err).Error("error in getting all products from the inventory")
		return nil, &common.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}
	}
	var products []ProductResp

	for _, p := range resp {
		unitPrice, _ := strconv.ParseFloat(p.UnitPrice, 64)
		products = append(products, ProductResp{
			Name:         p.Name,
			Quantity:     p.Quantity,
			Category:     p.Category,
			UnitPrice:    unitPrice,
			DateAdded:    p.DateAdded.Local().String(),
			DateModified: p.DateModified.Local().String(),
		})
	}

	return &GetAllProductsResp{
		Products: products,
	}, nil

}
