package analytics

import (
	"context"
	"database/sql"
	"net/http"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

func (v *Validator) GetMostOrdersPlaced(ctx context.Context, limit int) (*MostOrdersPlacedResp, *common.ErrorResponse) {
	return v.AnalyticsService.GetMostOrdersPlaced(ctx, limit)
}

func (a *AnalyticsServiceImpl) GetMostOrdersPlaced(ctx context.Context, limit int) (*MostOrdersPlacedResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetMostOrdersPlaced", "request": limit})

	//return rows based on the limit
	var ordersPlaced []OrdersPlaced

	resp, err := a.Queries.GetMostOrdersPlaced(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithError(err).Error("records not found. add orders first")
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "records not found. add orders first",
			}
		}
		logger.WithError(err).Error("error in get the most orders placed")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	l := 0
	for _, o := range resp {
		if l < limit {
			ordersPlaced = append(ordersPlaced, OrdersPlaced{
				Username:   o.Username,
				Email:      o.Email,
				OrderCount: int(o.OrdersCount),
			})
		} else {
			break
		}
		l++
	}

	return &MostOrdersPlacedResp{
		OrdersPlaced: ordersPlaced,
	}, nil
}
