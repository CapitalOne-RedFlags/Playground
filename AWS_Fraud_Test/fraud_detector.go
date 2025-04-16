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
// CreateModel creates a model if it doesn't exist
func CreateModel(client *frauddetector.Client) {
	modelID := "fraud_model_v2"

	// First check if model already exists
	existingModels, err := client.GetModels(context.TODO(), &frauddetector.GetModelsInput{})
	if err != nil {
		log.Printf("Warning: Unable to check existing models: %v", err)
	} else {
		for _, model := range existingModels.Models {
			if *model.ModelId == modelID {
				fmt.Printf("Model %s already exists, skipping creation.\n", modelID)
				return
			}
		}
	}

	// Create the model since it doesn't exist
	_, err = client.CreateModel(context.TODO(), &frauddetector.CreateModelInput{
		ModelId:       aws.String(modelID),
		ModelType:     types.ModelTypeEnumTransactionFraudInsights,
		EventTypeName: aws.String("transaction_event"),
		Description:   aws.String("Fraud detection model for transactions with enhanced variables"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
	fmt.Println("Model created successfully.")
}
