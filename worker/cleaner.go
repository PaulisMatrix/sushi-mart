package worker

import (
	"sushi-mart/common"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/sirupsen/logrus"
)

var CleanerLogger *logrus.Logger

// cleaner cleans up the consumer/queues keys and puts back all the unack entries to ready queue in case consuming is stopped for some reason

func Clean(logger *logrus.Logger) {
	// init logger
	CleanerLogger = logger
	CleanerLogger.Info("started the cleaner")

	// get the rmq connection
	conn := common.GetNewRMQConn(0)

	cleaner := rmq.NewCleaner(conn)
	// clean for every 1 min
	for range time.Tick(time.Minute) {
		returned, err := cleaner.Clean()
		if err != nil {
			CleanerLogger.WithError(err).Error("failed to clean")
			continue
		}
		CleanerLogger.WithField("num_cleaned", returned).Info("number of records cleaned")
	}

}
