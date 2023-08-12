package main

import (
	"net/http"
	"sushi-mart/api/analytics"
	"sushi-mart/api/inventory"
	"sushi-mart/api/user"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"sushi-mart/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func healthCheck(c *gin.Context) {
	logger := common.ExtractLoggerUnsafe(c.Request.Context())
	logger.WithField("method", "healthcheck").Info("healthcheck called")

	c.JSON(http.StatusOK, gin.H{"status": "pong!!"})
}

func setupMiddlewares(router *gin.RouterGroup, config *common.Config, logger *logrus.Logger) {
	//setup middlewares here
	router.Use(middlewares.LoggerMiddleware(logger))
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.JwtMiddleware(config))
}

func setupRoutes(engine *gin.Engine, Queries *database.Queries, config *common.Config, logger *logrus.Logger) {
	router := engine.Group("/api/v1")

	//add your routes here

	//routergroup for managing inventory, restricted to admins
	inventoryRouterGrp := router.Group("/inventory", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	inventory.New(Queries).HandleInventory(inventoryRouterGrp)

	//routergroup to check users,orders,products analytics, restricted to admins
	analyticsRouterGrp := router.Group("/analytics", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	analytics.New(Queries).HandleAnalytics(analyticsRouterGrp)

	//routergroup for managing users
	users := router.Group("/users")
	user.New(Queries).HandleUsers(users, config)

	//setup middlewares
	setupMiddlewares(router, config, logger)

	//jwt authenticated routes
	router.GET("/ping", healthCheck)
}
