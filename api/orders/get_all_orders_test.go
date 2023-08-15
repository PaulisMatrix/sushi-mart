package orders_test

import (
	"database/sql"
	"strconv"
	"sushi-mart/api/orders"
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
)

func (o *OrdersServiceSuite) TestGetAllOrdersOkReq() {
	Id := 1
	expectedResult := []database.GetAllPlacedOrdersRow{
		database.GetAllPlacedOrdersRow{
			OrderID:     1,
			OrderDate:   time.Now().Local(),
			OrderStatus: "PROCESSING",
			TotalAmt:    "1243.122",
			Username:    "testing",
			ProductName: "testProduct",
		},
	}

	var actualResult []orders.GetAllOrders

	for _, r := range expectedResult {
		amt, _ := strconv.ParseFloat(r.TotalAmt, 64)
		actualResult = append(actualResult, orders.GetAllOrders{
			OrderDate:   r.OrderDate.String(),
			OrderStatus: r.OrderStatus,
			TotalAmount: amt,
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
	expectedError := sql.ErrNoRows

	o.queriesMock.EXPECT().GetAllPlacedOrders(gomock.Any(), int32(Id)).Return(expectedResult, expectedError)
	resp, errResp := o.ordersService.GetOrders(o.context, Id)
	o.NotNil(errResp)
	o.EqualValues("no orders found", errResp.Message)
	o.Nil(resp)
}
