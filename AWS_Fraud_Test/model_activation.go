package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// ActivateModel deploys and activates the model if it's in TRAINING_COMPLETE status.
func ActivateModel(client *frauddetector.Client) {
	modelID := "fraud_model"
	modelVersion := "1.0"

	// Check model version status
	statusResp, err := client.GetModelVersion(context.TODO(), &frauddetector.GetModelVersionInput{
		ModelId:            aws.String(modelID),
		ModelVersionNumber: aws.String(modelVersion),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
	})
	if err != nil {
		log.Fatalf("Failed to get model version: %v", err)
	}

	// Print out the status to understand what you're receiving
	fmt.Printf("Model status: %v\n", statusResp.Status)

	// Only activate if status is TRAINING_COMPLETE
	if statusResp.Status == aws.String("TRAINING_COMPLETE") {
		_, err = client.UpdateModelVersionStatus(context.TODO(), &frauddetector.UpdateModelVersionStatusInput{
			ModelId:            aws.String(modelID),
			ModelType:          types.ModelTypeEnumTransactionFraudInsights,
			ModelVersionNumber: aws.String(modelVersion),
			Status:             types.ModelVersionStatusActive,
		})
		if err != nil {
			log.Fatalf("Failed to activate model: %v", err)
		}
		fmt.Println("Model activated successfully.")
	} else {
		fmt.Printf("Model is not in 'TRAINING_COMPLETE' status, current status: %v\n", statusResp.Status)
	}
}
