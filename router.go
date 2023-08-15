package main

import (
	"net/http"
	"sushi-mart/api/analytics"
	"sushi-mart/api/inventory"
	"sushi-mart/api/orders"
	"sushi-mart/api/user"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"sushi-mart/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func healthCheck(c *gin.Context) {
	logger := common.ExtractLoggerUnsafe(c.Request.Context())
	logger.WithField("method", "healthcheck").Info("healthcheck called")

	c.JSON(http.StatusOK, gin.H{"status": "pong!!"})
}

func helloAdmin(c *gin.Context) {
	logger := common.ExtractLoggerUnsafe(c.Request.Context())
	logger.WithError(nil).WithField("method", "helloAdmin").Error(" called")
	c.JSON(http.StatusOK, gin.H{"status": "granted admin access!!"})
}

/*
func setupMiddlewares(router *gin.RouterGroup, config *common.Config, logger *logrus.Logger) {
	//setup middlewares here
	router.Use(middlewares.LoggerMiddleware(logger))
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.JwtMiddleware(config))
}
*/

func setupRoutes(engine *gin.Engine, Queries *database.Queries, config *common.Config, logger *logrus.Logger) {
	router := engine.Group("/api/v1")

	//default middlewares
	router.Use(middlewares.LoggerMiddleware(logger))

	//add your routes here
	router.GET("/ping", healthCheck)

	//swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//admin routes
	router.GET("/admin/login", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}), helloAdmin)

	//routergroup for managing inventory, restricted to admins
	inventoryRouterGrp := router.Group("/admin/inventory", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	inventory.New(Queries).HandleInventory(inventoryRouterGrp)

	//routergroup to check users,orders,products analytics, restricted to admins
	analyticsRouterGrp := router.Group("/admin/analytics", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	analytics.New(Queries).HandleAnalytics(analyticsRouterGrp)

	//routergroup for managing users
	users := router.Group("/users")
	user.New(Queries).HandleUsers(users, config)

	//jwt authenticated routes
	router.Use(middlewares.JwtMiddleware(config))
	orderRouterGrp := router.Group("/orders")
	orders.New(Queries).HandleOrders(orderRouterGrp, config)
}
