package inventory

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (v *Validator) UpdateProduct(ctx context.Context, Id int, req *UpdateProductReq) (*ProductResp, *common.ErrorResponse) {
	return v.InventoryService.UpdateProduct(ctx, Id, req)
}

func (i *InventoryServiceImpl) UpdateProduct(ctx context.Context, Id int, req *UpdateProductReq) (*ProductResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "UpdateProduct", "request": req})

	dbParams := database.UpdateProductParams{ID: int32(Id), UpdateDateModified: true, DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true}}

	if req.Name == "" {
		dbParams.UpdateName = false
		dbParams.Name = req.Name
	} else {
		dbParams.UpdateName = true
		dbParams.Name = req.Name
	}

	if req.Quantity > 0 {
		dbParams.UpdateQuantity = true
		dbParams.Quantity = req.Quantity
	} else {
		dbParams.UpdateQuantity = false
		dbParams.Quantity = req.Quantity
	}

	if req.Category == "" {
		dbParams.UpdateCategory = false
		dbParams.Category = req.Category
	} else {
		dbParams.UpdateCategory = true
		dbParams.Category = req.Category
	}

	if req.UnitPrice > 0 {
		dbParams.UpdateUnitPrice = true
		dbParams.UnitPrice = decimal.NewFromFloat(req.UnitPrice)
	} else {
		dbParams.UpdateUnitPrice = false
		dbParams.UnitPrice = decimal.NewFromFloat(req.UnitPrice)
	}

	resp, err := i.Queries.UpdateProduct(ctx, dbParams)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithError(err).Error("record not found to update")
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "invalid id, record not found to update",
			}
		}

		logger.WithError(err).Error("failed to update the product in the inventory")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	return &ProductResp{
		Name:         resp.Name,
		Quantity:     resp.Quantity,
		Category:     resp.Category,
		UnitPrice:    resp.UnitPrice.Abs().InexactFloat64(),
		DateAdded:    resp.DateAdded.Time.Local().String(),
		DateModified: resp.DateModified.Time.Local().String(),
	}, nil
}
