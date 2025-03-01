package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Copied from SQS_Test

// NewSQSClient initializes and returns an SQS client.
func NewSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}
	return sqs.NewFromConfig(cfg)
}
