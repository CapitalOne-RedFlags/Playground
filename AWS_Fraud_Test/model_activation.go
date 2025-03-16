package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// ActivateModel deploys and activates the model.
func ActivateModel(client *frauddetector.Client) {
	modelID := "fraud_model"
	modelVersion := "1.0"

	_, err := client.UpdateModelVersionStatus(context.TODO(), &frauddetector.UpdateModelVersionStatusInput{
		ModelId:            aws.String(modelID),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
		ModelVersionNumber: aws.String(modelVersion),
		Status:             types.ModelVersionStatusActive,
	})
	if err != nil {
		log.Fatalf("Failed to activate model: %v", err)
	}
	fmt.Println("Model activated successfully.")
}
