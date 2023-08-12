package inventory

import (
	"context"
	"database/sql"
	"net/http"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

func (v *Validator) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {
	return v.InventoryService.DeleteProduct(ctx, id)
}

func (i *InventoryServiceImpl) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "DeleteProduct", "request": id})

	err := i.Queries.DeletProduct(ctx, int32(id))

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithError(err).Error("no rows matched to delete")
			return &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "invalid id, record not found",
			}
		}
		logger.WithError(err).Error("failed to delete the product from the inventory")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
