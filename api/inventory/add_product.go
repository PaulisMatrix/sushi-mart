package inventory

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (v *Validator) AddProduct(ctx context.Context, req *AddProductReq) *common.ErrorResponse {
	return v.InventoryService.AddProduct(ctx, req)
}

func (i *InventoryServiceImpl) AddProduct(ctx context.Context, req *AddProductReq) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "AddProduct", "request": req})

	dbParams := database.AddProductParams{
		Name:         req.Name,
		Quantity:     req.Quantity,
		Category:     req.Category,
		UnitPrice:    decimal.NewFromFloat(req.UnitPrice),
		DateAdded:    pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
	}
	err := i.Queries.AddProduct(ctx, dbParams)

	if err != nil {
		logger.WithError(err).Error("error in adding a new product to the inventory")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
