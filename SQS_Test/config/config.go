package config

import "os"

// GetQueueName retrieves the queue name from environment variables.
func GetQueueName() string {
	queueName := os.Getenv("SQS_QUEUE_NAME")
	if queueName == "" {
		queueName = "MyTestQueue" // Default if not set
	}
	return queueName
}
