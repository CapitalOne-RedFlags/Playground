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

	_, err := client.CreateModelVersion(context.TODO(), &frauddetector.CreateModelVersionInput{
		ModelId:            aws.String(modelID),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
		TrainingDataSource: types.TrainingDataSourceEnumIngestedEvents,
		TrainingDataSchema: &types.TrainingDataSchema{
			LabelSchema: &types.LabelSchema{
				LabelMapper: map[string][]string{
					"FRAUD": {"1"},
					"LEGIT": {"0"},
				},
				UnlabeledEventsTreatment: types.UnlabeledEventsTreatmentIgnore, // Move inside LabelSchema
			},
			ModelVariables: []string{"ip_address", "email_address", "transaction_amount"},
		},
		// UnlabeledEventsTreatmentEnum: types.UnlabeledEventsTreatmentIgnore, // ✅ Required when using INGESTED_EVENTS
	})
	if err != nil {
		log.Fatalf("Failed to train model: %v", err)
	}
	fmt.Println("Model training started.")
}
