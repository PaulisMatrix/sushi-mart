package main

//"net/http"

//"github.com/gin-gonic/gin"

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var DefaultLogger *logrus.Logger

// init logger
func init() {
	DefaultLogger = logrus.StandardLogger()
	DefaultLogger.SetFormatter(&logrus.JSONFormatter{})
}

func main() {

	//init db
	config := common.GetConfig()
	postgres, err := database.NewPostgres(config.PgDbName, config.PgUser, config.PgPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	queries := database.New(postgres.DB)

	r := gin.New()
	r.Use(gin.Recovery())

	//setup all the routes
	setupRoutes(r, queries, config, DefaultLogger)

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
