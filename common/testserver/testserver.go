package testserver

import (
	"net/http/httptest"
	"sushi-mart/common"
	"sushi-mart/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinServer struct {
	*httptest.Server
	Config      *common.Config
	RouterGroup *gin.RouterGroup
}

func GetNewTestGinServer() *GinServer {
	r := gin.New()
	r.Use(gin.Recovery())
	router := r.Group("/api/v1")

	// setup default test middleware
	defaultLogger := logrus.StandardLogger()
	defaultLogger.SetFormatter(&logrus.JSONFormatter{})
	setupDefaultTestMiddlewares(router, defaultLogger)

	s := httptest.NewUnstartedServer(r)
	return &GinServer{
		Server:      s,
		RouterGroup: router,
	}
}

func (s *GinServer) StartGinServer(prefix string) {
	s.Server.Start()
	s.URL = s.URL + prefix
}

func setupDefaultTestMiddlewares(router *gin.RouterGroup, logger *logrus.Logger) {
	router.Use(middlewares.LoggerMiddleware(logger))
}
