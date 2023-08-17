package orders_test

import (
	"errors"
	"sushi-mart/api/orders"
	"sushi-mart/internal/database"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (o *OrdersServiceSuite) TestPlaceOrderOkReq() {
	Id := 1
	req := &orders.PlaceOrderReq{
		Units:       5,
		ProductId:   1,
		PaymentType: "GPAY",
	}

	expectedProducItem := database.Productitem{
		ID:        1,
		Name:      "testing",
		Quantity:  100,
		Category:  "testingCategory",
		UnitPrice: decimal.NewFromFloat(34.21),
		IsActive:  true,
	}
	unitPrice := expectedProducItem.UnitPrice.Abs().InexactFloat64()
	totalAmt := unitPrice * float64(req.Units)

	// comment out OrderDate in PlaceOrder otherwise expected Date would be diff than actual Date due to change in time when mocking PlaceOrder
	dbParams := database.PlaceOrderParams{
		OrderStatus: string(orders.PROCESSING),
		TotalAmt:    decimal.NewFromFloat(totalAmt),
		Units:       int32(req.Units),
		PaymentType: req.PaymentType,
		CustomerID:  pgtype.Int4{Int32: int32(Id), Valid: true},
		ProductID:   pgtype.Int4{Int32: int32(req.ProductId), Valid: true},
	}

	o.queriesMock.EXPECT().GetProductItem(gomock.Any(), int32(req.ProductId)).Return(expectedProducItem, nil)
	o.queriesMock.EXPECT().PlaceOrder(gomock.Any(), dbParams).Return(nil)
	respErr := o.ordersService.PlaceOrder(o.context, req, Id)
	o.Nil(respErr)
}

func (o *OrdersServiceSuite) TestPlaceOrderBadReq() {
	Id := 2
	req := &orders.PlaceOrderReq{
		Units:       100,
		ProductId:   2,
		PaymentType: "GPAY",
	}

	expectedProducItem := database.Productitem{
		ID:        1,
		Name:      "testing",
		Quantity:  100,
		Category:  "testingCategory",
		UnitPrice: decimal.NewFromFloat(34.21),
		IsActive:  true,
	}
	unitPrice := expectedProducItem.UnitPrice.Abs().InexactFloat64()
	totalAmt := unitPrice * float64(req.Units)

	// comment out OrderDate in PlaceOrder otherwise expected Date would be diff than actual Date due to change in time when mocking PlaceOrder
	dbParams := database.PlaceOrderParams{
		OrderStatus: string(orders.PROCESSING),
		TotalAmt:    decimal.NewFromFloat(totalAmt),
		Units:       int32(req.Units),
		PaymentType: req.PaymentType,
		CustomerID:  pgtype.Int4{Int32: int32(Id), Valid: true},
		ProductID:   pgtype.Int4{Int32: int32(req.ProductId), Valid: true},
	}

	// when balance < total_amt
	expectedErr := errors.New("trigger failed.")

	o.queriesMock.EXPECT().GetProductItem(gomock.Any(), int32(req.ProductId)).Return(expectedProducItem, nil)
	o.queriesMock.EXPECT().PlaceOrder(gomock.Any(), dbParams).Return(expectedErr)
	respErr := o.ordersService.PlaceOrder(o.context, req, Id)
	o.NotNil(respErr)
	o.EqualValues(expectedErr.Error(), respErr.Message)
}
