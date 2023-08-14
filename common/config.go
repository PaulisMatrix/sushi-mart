package common

import (
	"log"
	"os"

	"github.com/adjust/rmq/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	PgDbName       string
	PgUser         string
	PgPass         string
	JwtSktKey      string
	AdminUser      string
	AdminPass      string
	QueueName      string
	RetryQueueName string
	OpenQueue      rmq.Queue
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil
	}
	return &Config{
		PgDbName:       os.Getenv("POSTGRES_DB"),
		PgUser:         os.Getenv("POSTGRES_USER"),
		PgPass:         os.Getenv("POSTGRES_PASS"),
		JwtSktKey:      os.Getenv("JWTSECRETKEY"),
		AdminUser:      os.Getenv("ADMIN_USER"),
		AdminPass:      os.Getenv("ADMIN_PASS"),
		QueueName:      os.Getenv("QUEUE_NAME"),
		RetryQueueName: os.Getenv("RETRY_QUEUE_NAME"),
		OpenQueue:      newRMQQueue(os.Getenv("QUEUE_NAME")),
	}
}

func GetNewRMQConn(errChanSize int) rmq.Connection {
	var errChan chan error
	if errChanSize == 0 {
		errChan = nil
	} else {
		errChan = make(chan error, errChanSize)
	}

	connection, err := rmq.OpenConnection("sushi-mart", "tcp", "localhost:6379", 1, errChan)
	if err != nil {
		log.Fatal("failed to establish a producer connection")
	}
	return connection
}

func newRMQQueue(queueName string) rmq.Queue {
	connection := GetNewRMQConn(0)
	newQueue, err := connection.OpenQueue(queueName)
	if err != nil {
		log.Fatal("failed to open a new queue")
	}

	return newQueue
}
