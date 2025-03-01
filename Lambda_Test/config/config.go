package config

import "os"

func GetQueueName() string {
	queueName := os.Getenv("SQS_QUEUE_NAME")
	if queueName == "" {
		queueName = "TransactionQueue" // Default if not set
	}
	return queueName
}
