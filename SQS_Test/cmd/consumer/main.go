package main

import (
	"log"

	"SQS_Test/config"
	"SQS_Test/pkg/sqs"
)

func main() {
	queueName := config.GetQueueName()

	// Ensure queue exists before consuming messages
	queueURL, err := sqs.GetQueueURL(queueName)
	if err != nil {
		log.Fatalf("Error retrieving queue URL: %v", err)
	}

	log.Println("Using Queue:", queueURL)

	sqsClient := sqs.NewSQSClient()
	sqs.ReceiveMessages(sqsClient, queueURL)
}
