package inventory

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/sirupsen/logrus"
)

func (v *Validator) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {
	return v.InventoryService.DeleteProduct(ctx, id)
}

func (i *InventoryServiceImpl) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "DeleteProduct", "request": id})

	dbParams := database.DeletProductParams{
		ID:       int32(id),
		IsActive: false,
	}

	numRows, err := i.Queries.DeletProduct(ctx, dbParams)

	if err != nil {
		logger.WithError(err).Error("failed to delete the product from the inventory")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	if numRows != 1 {
		logger.WithError(err).Error("no rows matched to delete")
		return &common.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "invalid id, record not found",
		}
	}

	return nil
}
