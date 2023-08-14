package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sushi-mart/api/orders"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"syscall"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/sirupsen/logrus"
)

var ConsumerLogger *logrus.Logger

const (
	prefetchLimit = 10
	pollDuration  = 100 * time.Millisecond

	retryFetchLimit   = 15
	retryPollDuration = 1 * time.Minute

	numConsumers = 3
	retryCount   = 3
)

type Consumer struct {
	name string
	orders.OrderService
}

func NewConsumer(queries *database.Queries, name string) *Consumer {
	return &Consumer{
		OrderService: &orders.OrderServiceImpl{
			Queries: queries,
		},
		name: name,
	}
}

func Consume(queries *database.Queries, config *common.Config, logger *logrus.Logger) {
	// init logger
	ConsumerLogger = logger

	ConsumerLogger.Info("started the worker")

	chanSize := 10
	errChan := make(chan error, chanSize)
	go handleErrors(errChan)

	// get the rmq connection
	conn := common.GetNewRMQConn(chanSize)

	// open the queue from which we would be consuming
	queue, err := conn.OpenQueue(config.QueueName)
	if err != nil {
		ConsumerLogger.WithError(err).Error("failed to get the corresponding queue")
		os.Exit(1)
	}

	// retry queue
	retryQueue, err := conn.OpenQueue(config.RetryQueueName)
	if err != nil {
		ConsumerLogger.WithError(err).Error("failed to get the corresponding retry queue")
		os.Exit(1)
	}

	// setup a retry queue
	queue.SetPushQueue(retryQueue)

	// start consuming
	// prefetchLimit is the limit rmq will fetch the records at once.
	// Always keep this higher than number of consumers available otherwise some of your consumers will remain idle.
	// pollDuration is the duration for which rmq will keep polling the records.
	if err := queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		ConsumerLogger.WithError(err).Error("failed to start consuming for the queue")
		os.Exit(1)
	}

	// start retry consumer
	if err := retryQueue.StartConsuming(retryFetchLimit, retryPollDuration); err != nil {
		ConsumerLogger.WithError(err).Error("failed to start consuming for retry queue")
		os.Exit(1)
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer %d", i)
		c := NewConsumer(queries, name)

		// After adding a consumer, it executes the Consume method in a separate go-routine already
		if _, err := queue.AddConsumer(name, c); err != nil {
			ConsumerLogger.WithError(err).Error("faield to start the consumers")
			os.Exit(1)
		}
	}

	// have only one cusomer for retry queue
	_, err = retryQueue.AddConsumer("retry-consumer", NewConsumer(queries, "retry-consumer"))
	if err != nil {
		ConsumerLogger.WithError(err).Error("faield to start the consumers")
		os.Exit(1)
	}

	// start a cleaner which will place all unack deliveries back into ready queue
	go func(conn rmq.Connection) {
		cleaner := rmq.NewCleaner(conn)
		// clean for every 1 min
		for range time.Tick(time.Minute) {
			returned, err := cleaner.Clean()
			if err != nil {
				ConsumerLogger.WithError(err).Error("failed to clean")
				continue
			}
			ConsumerLogger.WithField("num_cleaned", returned).Info("number of records cleaned")
		}
	}(conn)

	//start the stats server
	go statsServer(conn)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals

	ConsumerLogger.Info("shutting down all the consumers")

	// wait for all consumers to shutdown
	<-conn.StopAllConsuming()

}

// implement the Consume method here.
func (c *Consumer) Consume(delivery rmq.Delivery) {
	// Consumer will called PlaceOrder

	// get the PlaceOrder request
	var placeOrderReq orders.PlaceOrderReq
	payload := delivery.Payload()

	if err := json.Unmarshal([]byte(payload), &placeOrderReq); err != nil {
		// bad request, reject delivery
		// payload gets moved from unack queue to rejected queue
		// need a method to remove those rejected deliveries
		delivery.Reject()
	}

	// place the order
	err := c.OrderService.PlaceOrder(prepareCtx(context.Background()), &placeOrderReq, placeOrderReq.CustomerID)
	if err != nil {
		// processing order failed
		if err.Message == "trigger failed." {
			// report the status back to the user
			ConsumerLogger.Info("insufficient balance or not enough product units available to purchase")
		} else {
			// add to a retry queue
			delivery.Push()
		}
	}

	// ack, reject, push have in-built retry mechanism already
	ackErr := delivery.Ack()
	if ackErr != nil {
		if errors.Is(ackErr, rmq.ErrorConsumingStopped) {
			// consuming stopped. cleaner will move such delivers back into ready queue when consumers are up
			ConsumerLogger.WithError(ackErr).Error("consuming stopped")
		}

		ConsumerLogger.WithError(ackErr).Error("consuming stopped for unkown reason")
	}

	// report the status back to the user.
	ConsumerLogger.Info("Order registered succesfully!!!!")
}

func handleErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {

		case *rmq.ConsumeError:
			// Prefetching into delivery channel stopped for some reason.
			// but its fine since consumers will be idle for a moment then will again start consuming when new deliveres are fetched.
			ConsumerLogger.WithError(err).Error("consumer error")

		case *rmq.DeliveryError:
			// delivery error on ack, reject, push.
			// maintain a count to reject after X retries
			// this error lacks differentiating between whether its cause of ack, reject or push
			if err.Count >= retryCount {
				// reject the request
				ConsumerLogger.WithError(err).WithField("retry_count", retryCount).Error("requested rejected")
				err.Delivery.Reject()
			}
			ConsumerLogger.WithError(err).Error("delivery error")

		default:
			ConsumerLogger.WithError(err).Error("unkown error")
		}
	}
}

func prepareCtx(parentCtx context.Context) context.Context {
	updatedCtx := context.WithValue(parentCtx, common.LoggerKey{}, ConsumerLogger)
	return updatedCtx
}
