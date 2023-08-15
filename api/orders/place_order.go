package orders

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/sirupsen/logrus"
)

func (v *Validator) PlaceOrder(ctx context.Context, req *PlaceOrderReq, Id int) *common.ErrorResponse {
	return v.OrderService.PlaceOrder(ctx, req, Id)
}

func (o *OrderServiceImpl) PlaceOrder(ctx context.Context, req *PlaceOrderReq, Id int) *common.ErrorResponse {
	//before placing an order, calculate total_amt for the order
	//fetch unit_price from productItems
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "PlaceOrder", "request": req})

	prodResp, err := o.Queries.GetProductItem(ctx, int32(req.ProductId))
	if err != nil {
		logger.WithError(err).Error("failed to get the order product")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	unitPrice, _ := strconv.ParseFloat(prodResp.UnitPrice, 64)
	totalAmt := unitPrice * float64(req.Units)

	//there are two triggers
	//one for validating totalAmt against balance from wallet table
	//one for updating quantity in productItems
	dbParams := database.PlaceOrderParams{
		OrderStatus: string(PROCESSING),
		TotalAmt:    strconv.FormatFloat(totalAmt, 'E', -1, 64),
		Units:       int32(req.Units),
		PaymentType: req.PaymentType,
		OrderDate:   time.Now().Local(),
		CustomerID:  sql.NullInt32{Int32: int32(Id), Valid: true},
		ProductID:   sql.NullInt32{Int32: int32(req.ProductId), Valid: true},
	}
	orderErr := o.Queries.PlaceOrder(ctx, dbParams)

	// use pgx driver here to get explicit error
	if orderErr != nil {
		//trigger failed to run
		logger.WithError(orderErr).Error("failed to place the order")
		return &common.ErrorResponse{
			Status:  http.StatusOK,
			Message: "trigger failed.",
		}
	}

	return nil
}
