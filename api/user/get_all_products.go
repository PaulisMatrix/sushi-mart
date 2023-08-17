package user

import (
	"context"
	"net/http"
	"sushi-mart/common"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (v *Validator) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	return v.UsersService.GetAllProducts(ctx)
}

func (u *UsersServiceImpl) GetAllProducts(ctx context.Context) (*GetAllProductsResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetAllProducts"})

	resp, err := u.Queries.GetAllProducts(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithError(err).Error("no products in the inventory")
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "no products in the inventory",
			}
		}
		logger.WithError(err).Error("error in getting different products")
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
