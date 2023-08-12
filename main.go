package main

//"net/http"

//"github.com/gin-gonic/gin"

import (
	"fmt"
	"log"
	"os"
	"sushi-mart/common"
	"sushi-mart/internal/database"

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

	switch os.Args[1] {
	case "serve":
		//start the server
		server(queries, config)
	case "bg-worker":
		//start the bg worker
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}

}
