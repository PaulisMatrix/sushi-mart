package orders

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

func (v *Validator) GetOrders(ctx context.Context, Id int) (*GetAllOrdersResp, *common.ErrorResponse) {
	return v.OrderService.GetOrders(ctx, Id)
}

func (o *OrderServiceImpl) GetOrders(ctx context.Context, Id int) (*GetAllOrdersResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetOrders", "request": Id})

	var orders []GetAllOrders

	resp, err := o.Queries.GetAllPlacedOrders(ctx, int32(Id))

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithError(err).Info("no orders found for given customer id")
			return nil, &common.ErrorResponse{
				Status:  http.StatusOK,
				Message: "no orders found",
			}
		}
		logger.WithError(err).Error("failed to get the orders for the given customer")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	for _, o := range resp {
		amt, _ := strconv.ParseFloat(o.TotalAmt, 64)
		orders = append(orders, GetAllOrders{
			OrderDate:   o.OrderDate.String(),
			OrderStatus: o.OrderStatus,
			TotalAmount: amt,
			Username:    o.Username,
			ProductName: o.ProductName,
		})
	}
	return &GetAllOrdersResp{
		Orders: orders,
	}, nil
}
