package main

import (
	"log"

	"SQS_Test/config"
	"SQS_Test/pkg/sqs"
)

func main() {
	queueName := config.GetQueueName()

	// Step 1: Ensure the queue exists
	queueURL, err := sqs.CreateQueue(queueName)
	if err != nil {
		log.Fatalf("Error creating queue: %v", err)
	}

	log.Println("Using Queue:", queueURL)

	sqsClient := sqs.NewSQSClient()

	// Send multiple transactions
	transactions := []sqs.Transaction{
		{TransactionID: "txn12345", Amount: 100.00},
		{TransactionID: "txn67890", Amount: 250.00},
		{TransactionID: "txn54321", Amount: 75.50},
	}

	for _, txn := range transactions {
		err := sqs.SendTransaction(sqsClient, txn, queueURL)
		if err != nil {
			log.Printf("Error sending transaction %s: %v", txn.TransactionID, err)
		}
	}
}
