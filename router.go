package main

import (
	"net/http"
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
	//setup jwt middleware here
	router.Use(middlewares.JwtMiddleware(config))
}

func setupRoutes(engine *gin.Engine, Queries *database.Queries, config *common.Config) {
	router := engine.Group("/api/v1")
	//add your routes here

	// routergroup for managing users
	users := router.Group("/users")
	user.New(Queries).HandlUsers(users, config)

	//setup middlewares
	setupMiddlewares(router, config)

	//routergroup for managing inventory
	//inventory := router.Group("/inventory")

	router.GET("/ping", healthCheck)
}
