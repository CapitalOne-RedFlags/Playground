package sqs

import (
	"Lambda_Test/internal/domain"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events" // For prototype only, events will be moved from sqs
)

func ProcessRecords(messages []events.SQSMessage) []domain.Transaction {
	var transactions []domain.Transaction
	for _, msg := range messages {
		var txn domain.Transaction
		err := json.Unmarshal([]byte(msg.Body), &txn)
		if err != nil {
			fmt.Printf("Failed to parse message: %v", err)
			continue
		}

		fmt.Printf("Processing Transaction: ID=%s, Amount=%.2f\n", txn.Id, txn.Amount)
		transactions = append(transactions, txn)
	}

	return transactions
}
