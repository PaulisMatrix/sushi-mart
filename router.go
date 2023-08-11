package main

import (
	"net/http"
	"sushi-mart/api/inventory"
	"sushi-mart/api/user"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"sushi-mart/middlewares"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pong!!"})
}

func setupMiddlewares(router *gin.RouterGroup, config *common.Config) {
	//setup middlewares here
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.JwtMiddleware(config))
}

func setupRoutes(engine *gin.Engine, Queries *database.Queries, config *common.Config) {
	router := engine.Group("/api/v1")

	//add your routes here

	//routergroup for managing inventory, restricted to admins
	invent := router.Group("/inventory", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	inventory.New(Queries).HandleInventory(invent)

	//routergroup for managing users
	users := router.Group("/users")
	user.New(Queries).HandleUsers(users, config)

	//setup middlewares
	setupMiddlewares(router, config)

	//jwt authenticated routes
	router.GET("/ping", healthCheck)
}
