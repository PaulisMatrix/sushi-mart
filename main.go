package main

import (
	"fmt"
	"log"
	"os"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"sushi-mart/worker"

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
	defer postgres.DB.Close()

	switch os.Args[1] {
	case "serve":
		//start the server
		server(queries, config)
	case "consume":
		//start the bg consumer
		worker.Consume(queries, config, DefaultLogger)
	case "clean":
		//start the cleaner
		worker.Clean(DefaultLogger)
	case "test-mig":
		database.MigrationsTest()
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}

}
