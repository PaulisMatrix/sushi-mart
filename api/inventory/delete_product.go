package inventory

import (
	"context"
	"database/sql"
	"net/http"
	"sushi-mart/common"
)

func (v *Validator) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {
	return v.InventoryService.DeleteProduct(ctx, id)
}

func (i *InventoryServiceImpl) DeleteProduct(ctx context.Context, id int) *common.ErrorResponse {

	err := i.Queries.DeletProduct(ctx, int32(id))

	if err != nil {
		if err == sql.ErrNoRows {
			return &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "invalid id, record not found",
			}
		}
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
