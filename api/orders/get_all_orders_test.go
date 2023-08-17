package orders_test

import (
	"sushi-mart/api/orders"
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (o *OrdersServiceSuite) TestGetAllOrdersOkReq() {
	Id := 1
	expectedResult := []database.GetAllPlacedOrdersRow{
		{
			OrderID:     1,
			OrderDate:   pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
			OrderStatus: "PROCESSING",
			TotalAmt:    decimal.NewFromFloat(1243.122),
			Username:    "testing",
			ProductName: "testProduct",
		},
	}

	var actualResult []orders.GetAllOrders

	for _, r := range expectedResult {
		actualResult = append(actualResult, orders.GetAllOrders{
			OrderDate:   r.OrderDate.Time.Local().String(),
			OrderStatus: r.OrderStatus,
			TotalAmount: r.TotalAmt.Abs().InexactFloat64(),
			Username:    r.Username,
			ProductName: r.ProductName,
		})
	}

	o.queriesMock.EXPECT().GetAllPlacedOrders(gomock.Any(), int32(Id)).Return(expectedResult, nil)
	resp, errResp := o.ordersService.GetOrders(o.context, Id)
	o.Nil(errResp)
	o.NotNil(resp)
	o.EqualValues(actualResult[0].Username, expectedResult[0].Username)
}

func (o *OrdersServiceSuite) TestGetAllOrdersNoRows() {
	Id := 2
	expectedResult := []database.GetAllPlacedOrdersRow{}
	expectedError := pgx.ErrNoRows

	o.queriesMock.EXPECT().GetAllPlacedOrders(gomock.Any(), int32(Id)).Return(expectedResult, expectedError)
	resp, errResp := o.ordersService.GetOrders(o.context, Id)
	o.NotNil(errResp)
	o.EqualValues(expectedError.Error(), "no rows in result set")
	o.Nil(resp)
}
