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
)

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
	setupRoutes(r, queries, config)

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
	/*
		password := "testing"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("ERROR IN GENERATING PASSWORD")
		}

		//create a user
		insertedUser, err := queries.CreateCustomer(ctx, database.CreateCustomerParams{
			Username: "rushikesh yadwade",
			Password: string(hashedPassword),
			Email:    "rushi@gmail.com",
			Phone:    sql.NullString{String: "9834150521", Valid: true},
			Address:  sql.NullString{String: "Miraj,Maharahstra", Valid: true},
		})

		if err != nil {
			fmt.Println("ERROR IN INSERTING", err)
		}

		fmt.Println("INSERTED USER", insertedUser)

		fetchedAuthor, err := queries.GetCustomer(ctx, "rushi@gmail.com")
		fmt.Println(fetchedAuthor)

		err = bcrypt.CompareHashAndPassword([]byte(fetchedAuthor.Password), []byte(password))

		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			fmt.Println("WRONG PASSWORD")
		}
	*/
}
