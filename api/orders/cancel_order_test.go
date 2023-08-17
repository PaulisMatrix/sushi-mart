package orders_test

import (
	"sushi-mart/api/orders"
	"sushi-mart/internal/database"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
)

func (o *OrdersServiceSuite) TestCancelOrderOkReq() {
	req := &orders.UpdateOrderReq{
		OrderId: 1,
	}
	dbParams := database.CancelOrderParams{
		ID:            int32(req.OrderId),
		OrderStatus:   string(orders.PROCESSING),
		OrderStatus_2: string(orders.CANCELLED),
	}

	var expectedNumRows int64
	expectedNumRows = 1
	o.queriesMock.EXPECT().CancelOrder(gomock.Any(), dbParams).Return(expectedNumRows, nil)
	errResp := o.ordersService.CancelOrder(o.context, req)
	o.Nil(errResp)
}

func (o *OrdersServiceSuite) TestCancelOrderBadReq() {
	req := &orders.UpdateOrderReq{
		OrderId: 1,
	}
	dbParams := database.CancelOrderParams{
		ID:            int32(req.OrderId),
		OrderStatus:   string(orders.PROCESSING),
		OrderStatus_2: string(orders.CANCELLED),
	}

	var expectedNumRows int64
	expectedNumRows = 0
	expectedMsg := "cannot cancel a shipped/delivered order"
	o.queriesMock.EXPECT().CancelOrder(gomock.Any(), dbParams).Return(expectedNumRows, nil)
	errResp := o.ordersService.CancelOrder(o.context, req)
	o.NotNil(errResp)
	o.EqualValues(expectedMsg, errResp.Message)
}

func (o *OrdersServiceSuite) TestCancelOrderNoRowsReq() {
	req := &orders.UpdateOrderReq{
		OrderId: 1,
	}
	dbParams := database.CancelOrderParams{
		ID:            int32(req.OrderId),
		OrderStatus:   string(orders.PROCESSING),
		OrderStatus_2: string(orders.CANCELLED),
	}

	var expectedNumRows int64
	expectedNumRows = 0
	expectedError := pgx.ErrNoRows
	o.queriesMock.EXPECT().CancelOrder(gomock.Any(), dbParams).Return(expectedNumRows, expectedError)
	errResp := o.ordersService.CancelOrder(o.context, req)
	o.NotNil(errResp)
	o.EqualValues(expectedError.Error(), "no rows in result set")
}
