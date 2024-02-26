package tests

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/common/testserver"
	mock_queries "sushi-mart/internal/database/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueryParams struct {
	key   string
	value string
}

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
	s.server = makeTestServerClient()
	s.client = s.server.Client()
}

func (s *ServerSuite) TearDownSuite() {
	s.queriesMock = nil
	s.server.Close()
	s.ctrl.Finish()
}

func makeTestServerClient() *testserver.GinServer {
	server := testserver.GetNewTestGinServer()

	// setup a testing route
	SetupPingRoute(server.RouterGroup)

	// start the server
	prefix := server.RouterGroup.BasePath()
	server.StartGinServer(prefix)

	return server
}

func SetupPingRoute(router *gin.RouterGroup) {
	router.Group("/testing").GET("/ping", healthCheck)
}

func healthCheck(c *gin.Context) {
	logger := common.ExtractLoggerUnsafe(c.Request.Context())
	logger.WithField("method", "healthcheck").Info("healthcheck called")

	c.JSON(http.StatusOK, gin.H{"status": "pong!!"})
}

func (s *ServerSuite) makeRequest(path string, params ...QueryParams) ([]byte, error) {
	req, _ := http.NewRequest("GET", s.server.URL+path, bytes.NewReader([]byte("")))
	newReq := req.URL.Query()
	for _, p := range params {
		newReq.Add(p.key, p.value)
	}
	req.URL.RawQuery = newReq.Encode()

	req.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return resBody, nil
}
