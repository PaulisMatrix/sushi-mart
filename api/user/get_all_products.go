package user

import (
	"context"
	"net/http"
	"strconv"
	"sushi-mart/common"
)

func (v *Validator) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	return v.UsersService.GetAllProducts(ctx)
}

func (u *UsersServiceImpl) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	resp, err := u.Queries.GetAllProducts(ctx)
	if err != nil {
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
