package orders

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/sirupsen/logrus"
)

func (v *Validator) CancelOrder(ctx context.Context, req *UpdateOrderReq) *common.ErrorResponse {
	return v.OrderService.CancelOrder(ctx, req)
}

func (o *OrderServiceImpl) CancelOrder(ctx context.Context, req *UpdateOrderReq) *common.ErrorResponse {
	//cancel order only if the status is PROCESSING
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "CancelOrder", "request": req})

	dbParams := database.UpdateOrderStatusParams{
		ID:            int32(req.OrderId),
		OrderStatus:   string(PROCESSING),
		OrderStatus_2: string(CANCELLED),
	}
	resp, err := o.Queries.UpdateOrderStatus(ctx, dbParams)

	if err != nil {
		logger.WithError(err).Error("failed to cancel the order")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	if resp != 1 {
		return &common.ErrorResponse{
			Status:  http.StatusOK,
			Message: "cannot cancel a shipped/delivered order",
		}
	}

	return nil
}
