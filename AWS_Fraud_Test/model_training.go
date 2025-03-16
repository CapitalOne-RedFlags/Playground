package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// TrainModel starts model training.
func TrainModel(client *frauddetector.Client) {
	modelID := "fraud_model"

	dataSource := types.TrainingDataSourceEnumExternalEvents

	_, err := client.CreateModelVersion(context.TODO(), &frauddetector.CreateModelVersionInput{
		ModelId:            aws.String(modelID),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
		TrainingDataSource: dataSource,
		TrainingDataSchema: &types.TrainingDataSchema{
			LabelSchema: &types.LabelSchema{
				LabelMapper: map[string][]string{
					"fraud": {"1"},
					"legit": {"0"},
				},
			},
			ModelVariables: []string{"ip_address", "email_address", "transaction_amount"},
		},
		// Add ExternalEventsDetail field which is required when using EXTERNAL_EVENTS
		ExternalEventsDetail: &types.ExternalEventsDetail{
			DataAccessRoleArn: aws.String("arn:aws:iam::YOUR_ACCOUNT_ID:role/service-role/AmazonFraudDetectorRole"),
			DataLocation:      aws.String("s3://redflags-bucket/bank_transactions_data.csv"),
			// EventType:         aws.String("transaction_event"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to train model: %v", err)
	}
	fmt.Println("Model training started.")
}
