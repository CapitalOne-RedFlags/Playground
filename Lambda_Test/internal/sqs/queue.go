package sqs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Copied from SQS_Test

// CreateQueue ensures the queue exists and returns its URL.
func CreateQueue(queueName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error loading AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	// Create the queue
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]string{
			"VisibilityTimeout": "30", // Message remains hidden for 30 seconds after being received
		},
	}

	result, err := sqsClient.CreateQueue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to create queue: %w", err)
	}

	log.Println("Queue Created:", *result.QueueUrl)
	return *result.QueueUrl, nil
}

// GetQueueURL retrieves the SQS queue URL
func GetQueueURL(queueName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error loading AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := sqsClient.GetQueueUrl(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get queue URL: %w", err)
	}

	log.Println("Queue URL Retrieved:", *result.QueueUrl)
	return *result.QueueUrl, nil
}
