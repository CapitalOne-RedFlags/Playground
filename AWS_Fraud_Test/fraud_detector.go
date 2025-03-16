package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// CreateModel registers a new fraud detection model.
func CreateModel(client *frauddetector.Client) {
	modelID := "fraud_model"
	description := "Fraud detection model"

	_, err := client.CreateModel(context.TODO(), &frauddetector.CreateModelInput{
		ModelId:       aws.String(modelID),
		ModelType:     types.ModelTypeEnumTransactionFraudInsights,
		Description:   aws.String(description),
		EventTypeName: aws.String("transaction_event"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
	fmt.Println("Model created successfully.")
}
