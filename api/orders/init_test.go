package orders_test

import (
	"context"
	"sushi-mart/api/orders"
	"sushi-mart/common"
	mock_queries "sushi-mart/internal/database/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrdersServiceSuite struct {
	suite.Suite
	*require.Assertions
	ctrl          *gomock.Controller
	queriesMock   *mock_queries.MockQuerier
	context       context.Context
	ordersService orders.OrderService
}

func TestOrdersServiceSuite(t *testing.T) {
	suite.Run(t, new(OrdersServiceSuite))
}

func (o *OrdersServiceSuite) SetupSuite() {
	o.Assertions = require.New(o.Suite.T())
	o.ctrl = gomock.NewController(o.T())
	o.queriesMock = mock_queries.NewMockQuerier(o.ctrl)
	o.context = prepareCtx(context.Background())
	o.ordersService = getOrdersService(o.queriesMock)
}

func (o *OrdersServiceSuite) TearDownSuite() {
	o.ctrl.Finish()
	o.queriesMock = nil
}

func prepareCtx(parentCtx context.Context) context.Context {
	updatedCtx := context.WithValue(parentCtx, common.LoggerKey{}, logrus.StandardLogger())
	return updatedCtx
}

func getOrdersService(mockQueries *mock_queries.MockQuerier) *orders.RoutesWrapper {
	return &orders.RoutesWrapper{
		OrderService: &orders.Validator{
			OrderService: &orders.Cache{
				OrderService: &orders.OrderServiceImpl{
					Queries: mockQueries,
				},
			},
		},
	}
}
