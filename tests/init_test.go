package tests

import (
	"context"
	"net/http"
	"sushi-mart/common/testserver"
	mock_queries "sushi-mart/internal/database/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	*require.Assertions
	ctrl        *gomock.Controller
	queriesMock *mock_queries.MockQuerier
	context     context.Context
	server      *testserver.GinServer
	client      *http.Client
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (s *ServerSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.ctrl = gomock.NewController(s.T())
	s.queriesMock = mock_queries.NewMockQuerier(s.ctrl)
}

func makeTestServerClient() (*testserver.GinServer, *http.Client) {
	return nil, nil
}
