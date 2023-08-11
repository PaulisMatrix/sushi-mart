package inventory

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"
)

func (v *Validator) UpdateProduct(ctx context.Context, Id int, req *UpdateProductReq) (*ProductResp, *common.ErrorResponse) {
	return v.InventoryService.UpdateProduct(ctx, Id, req)
}

func (i *InventoryServiceImpl) UpdateProduct(ctx context.Context, Id int, req *UpdateProductReq) (*ProductResp, *common.ErrorResponse) {

	dbParams := database.UpdateProductParams{ID: int32(Id), UpdateDateModified: true, DateModified: time.Now().Local()}

	if req.Name != "" {
		dbParams.UpdateName = true
		dbParams.Name = req.Name
	}

	if req.Quantity != 0 {
		dbParams.UpdateQuantity = true
		dbParams.Quantity = req.Quantity
	}

	if req.Category != "" {
		dbParams.UpdateCategory = true
		dbParams.Category = req.Category
	}

	if req.UnitPrice != 0 {
		dbParams.UpdateUnitPrice = true
		dbParams.UnitPrice = strconv.FormatFloat(req.UnitPrice, 'E', -1, 64)
	}

	resp, err := i.Queries.UpdateProduct(ctx, dbParams)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "invalid id, record not found to update",
			}
		}

		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	unitPrice, _ := strconv.ParseFloat(resp.UnitPrice, 64)

	return &ProductResp{
		Name:         resp.Name,
		Quantity:     resp.Quantity,
		Category:     resp.Category,
		UnitPrice:    unitPrice,
		DateAdded:    resp.DateAdded.Local().String(),
		DateModified: resp.DateModified.Local().String(),
	}, nil
}
