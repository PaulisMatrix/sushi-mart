package user_test

import (
	"context"
	"sushi-mart/api/user"
	"sushi-mart/common"
	mock_queries "sushi-mart/internal/database/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UsersServiceSuite struct {
	suite.Suite
	*require.Assertions
	ctrl         *gomock.Controller
	queriesMock  *mock_queries.MockQuerier
	context      context.Context
	UsersService user.UsersService
}

func TestOrdersServiceSuite(t *testing.T) {
	suite.Run(t, new(UsersServiceSuite))
}

func (u *UsersServiceSuite) SetupSuite() {
	u.Assertions = require.New(u.Suite.T())
	u.ctrl = gomock.NewController(u.T())
	u.queriesMock = mock_queries.NewMockQuerier(u.ctrl)
	u.context = prepareCtx(context.Background())
	u.UsersService = getUsersService(u.queriesMock)
}

func (u *UsersServiceSuite) TearDownSuite() {
	u.ctrl.Finish()
	u.queriesMock = nil
}

func prepareCtx(parentCtx context.Context) context.Context {
	updatedCtx := context.WithValue(parentCtx, common.LoggerKey{}, logrus.StandardLogger())
	return updatedCtx
}

func getUsersService(mockQueries *mock_queries.MockQuerier) *user.RoutesWrapper {
	return &user.RoutesWrapper{
		UsersService: &user.Validator{
			UsersService: &user.Cache{
				UsersService: &user.UsersServiceImpl{
					Queries: mockQueries,
				},
			},
		},
	}
}
