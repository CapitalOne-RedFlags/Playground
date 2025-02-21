package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Global queue URL variable
var queueURL string

// Create an SQS queue with a retry setting of 3
func createSQSQueue(queueName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	// Create the queue with attributes (max receive count = 3)
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]string{
			"VisibilityTimeout": "30",  // Message remains hidden for 30s after being received
			"RedrivePolicy": fmt.Sprintf(`{"maxReceiveCount":"3","deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:DeadLetterQueue"}`), // Retry 3 times before moving to DLQ
		},
	}

	result, err := sqsClient.CreateQueue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to create SQS queue: %v", err)
	}

	fmt.Println("SQS Queue Created:", *result.QueueUrl)
	return *result.QueueUrl, nil
}

// Send a transaction message to the SQS queue
func sendTransaction(sqsClient *sqs.Client, transactionID, amount string) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(fmt.Sprintf(`{"transaction_id":"%s","amount":"%s"}`, transactionID, amount)),
	}

	_, err := sqsClient.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	fmt.Println("Transaction sent to SQS:", transactionID)
	return nil
}

// Delete the SQS queue
func deleteSQSQueue(sqsClient *sqs.Client, queueURL string) error {
	input := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	}

	_, err := sqsClient.DeleteQueue(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete SQS queue: %v", err)
	}

	fmt.Println("SQS Queue Deleted:", queueURL)
	return nil
}

func main() {
	queueName := "TransactionQueue"

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	// Step 1: Create the Queue
	queueURL, err = createSQSQueue(queueName)
	if err != nil {
		log.Fatalf("Error creating SQS queue: %v", err)
	}

	// Step 2: Send Transactions
	err = sendTransaction(sqsClient, "txn12345", "100.00")
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	err = sendTransaction(sqsClient, "txn67890", "250.00")
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	// Wait to simulate message processing
	fmt.Println("Waiting for messages to be processed...")
	time.Sleep(10 * time.Second)

	// Step 3: Delete the Queue
	err = deleteSQSQueue(sqsClient, queueURL)
	if err != nil {
		log.Fatalf("Error deleting SQS queue: %v", err)
	}
}
