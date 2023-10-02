package common

import (
	"log"
	"os"

	"github.com/adjust/rmq/v5"
	"github.com/joho/godotenv"
)

var producerQueue rmq.Queue

var producerConn rmq.Connection

type Config struct {
	PgDbName       string
	PgTestDbName   string
	PgUser         string
	PgPass         string
	JwtSktKey      string
	AdminUser      string
	AdminPass      string
	QueueName      string
	RetryQueueName string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
		return nil
	}
	return &Config{
		PgDbName:       os.Getenv("POSTGRES_DB"),
		PgTestDbName:   os.Getenv("POSTGRES_TESTDB"),
		PgUser:         os.Getenv("POSTGRES_USER"),
		PgPass:         os.Getenv("POSTGRES_PASS"),
		JwtSktKey:      os.Getenv("JWTSECRETKEY"),
		AdminUser:      os.Getenv("ADMIN_USER"),
		AdminPass:      os.Getenv("ADMIN_PASS"),
		QueueName:      os.Getenv("QUEUE_NAME"),
		RetryQueueName: os.Getenv("RETRY_QUEUE_NAME"),
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

func getProducerConnection() rmq.Connection {
	if producerConn != nil {
		return producerConn
	}
	producerConn = GetNewRMQConn(0)
	return producerConn
}

func GetProducerQueue(queueName string) (rmq.Queue, error) {
	conn := getProducerConnection()

	if producerQueue != nil {
		return producerQueue, nil
	}

	producerQueue, err := conn.OpenQueue(queueName)
	if err != nil {
		log.Fatal("failed to open a new queue")
		return nil, err
	}
	return producerQueue, nil
}
