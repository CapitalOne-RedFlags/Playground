package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Transaction represents a transaction message.
type Transaction struct {
	TransactionID string `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

// SendTransaction sends a transaction message to the SQS queue.
func SendTransaction(sqsClient *sqs.Client, txn Transaction, queueURL string) error {
	body, err := json.Marshal(txn)
	if err != nil {
		return fmt.Errorf("error encoding transaction: %v", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	}

	_, err = sqsClient.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	log.Println("Transaction sent to SQS:", txn.TransactionID)
	return nil
}
