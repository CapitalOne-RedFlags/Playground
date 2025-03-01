package handlers

import (
	"Lambda_Test/internal/db"
	"Lambda_Test/internal/sqs"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func BatchProcessingHandler(ctx context.Context, event events.SQSEvent) {
	transactions := sqs.ProcessRecords(event.Records)

	for _, txn := range transactions {
		db.InsertTransaction(txn)
	}
}
