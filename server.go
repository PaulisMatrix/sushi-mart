package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sushi-mart/common"
	"sushi-mart/docs"
	"sushi-mart/internal/database"
	"sushi-mart/middlewares"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:generate $GOBIN/swag init --parseDependency --parseInternal -g ./server.go

// @title sushimart
// @description An OrderManagement service
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @description API Key to be provided for authentication
// @securityDefinitions.basic BasicAuth

func server(queries *database.Queries, config *common.Config) {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middlewares.CORSMiddleware())
	//setup all the routes
	setupRoutes(r, queries, config, DefaultLogger)
	docs.SwaggerInfo.BasePath = "/api/v1"

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("server exiting...")
	}
}
