package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// ReceiveMessages polls messages from SQS and processes them.
func ReceiveMessages(sqsClient *sqs.Client, queueURL string) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 5,
		WaitTimeSeconds:     10,
	}

	result, err := sqsClient.ReceiveMessage(context.TODO(), input)
	if err != nil {
		log.Printf("Error receiving messages: %v", err)
		return
	}

	for _, msg := range result.Messages {
		var txn Transaction
		err := json.Unmarshal([]byte(*msg.Body), &txn)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		fmt.Printf("Processing Transaction: ID=%s, Amount=%.2f\n", txn.TransactionID, txn.Amount)

		// Delete the message after processing
		deleteMessage(sqsClient, queueURL, *msg.ReceiptHandle)
	}
}

// deleteMessage removes a processed message from the queue.
func deleteMessage(sqsClient *sqs.Client, queueURL, receiptHandle string) {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := sqsClient.DeleteMessage(context.TODO(), input)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return
	}

	log.Println("Message deleted from SQS")
}
