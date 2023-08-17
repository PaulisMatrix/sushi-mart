package orders

import (
	"context"
	"errors"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
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
	unitPrice := prodResp.UnitPrice.Abs().InexactFloat64()
	totalAmt := unitPrice * float64(req.Units)

	//there are two triggers
	//one for validating totalAmt against balance from wallet table
	//one for updating quantity in productItems
	dbParams := database.PlaceOrderParams{
		OrderStatus: string(PROCESSING),
		TotalAmt:    decimal.NewFromFloat(totalAmt),
		Units:       int32(req.Units),
		PaymentType: req.PaymentType,
		OrderDate:   pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID:  pgtype.Int4{Int32: int32(Id), Valid: true},
		ProductID:   pgtype.Int4{Int32: int32(req.ProductId), Valid: true},
	}
	orderErr := o.Queries.PlaceOrder(ctx, dbParams)

	if orderErr != nil {
		//trigger failed to run
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
			logger.WithError(orderErr).Error("failed to place the order")
			return &common.ErrorResponse{
				Status:  http.StatusOK,
				Message: "trigger failed.",
			}
		} else {
			logger.WithError(err).Error("internal server eror")
			return &common.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			}
		}

	}

	return nil
}
