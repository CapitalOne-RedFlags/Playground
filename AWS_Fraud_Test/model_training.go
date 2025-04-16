package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// TrainModel starts model training and waits for it to complete.
func TrainModel(client *frauddetector.Client) {
	modelID := "fraud_model"

	output, err := client.CreateModelVersion(context.TODO(), &frauddetector.CreateModelVersionInput{
		ModelId:            aws.String(modelID),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
		TrainingDataSource: types.TrainingDataSourceEnumIngestedEvents,
		TrainingDataSchema: &types.TrainingDataSchema{
			LabelSchema: &types.LabelSchema{
				LabelMapper: map[string][]string{
					"FRAUD": {"FRAUD"},
					"LEGIT": {"LEGIT"},
				},
				UnlabeledEventsTreatment: types.UnlabeledEventsTreatmentIgnore,
			},
			ModelVariables: []string{"ip_address", "email_address", "transaction_amount"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to train model: %v", err)
	}

	version := aws.ToString(output.ModelVersionNumber)
	fmt.Printf("Model training started for version %s...\n", version)

	for {
		time.Sleep(30 * time.Second)

		statusResp, err := client.GetModelVersion(context.TODO(), &frauddetector.GetModelVersionInput{
			ModelId:            aws.String(modelID),
			ModelType:          types.ModelTypeEnumTransactionFraudInsights,
			ModelVersionNumber: aws.String(version),
		})
		if err != nil {
			log.Printf("Error polling model status: %v", err)
			continue
		}

		status := aws.ToString(statusResp.Status)
		fmt.Printf("Training status: %s\n", status)

		if status == "TRAINING_COMPLETE" {
			fmt.Println("✅ Training completed successfully.")
			break
		} else if status == "ERROR" {
			fmt.Println("❌ Training failed. Attempting to get detailed reason...")
			describeModelVersion(client, modelID, version)
			break
		}
	}
}

func describeModelVersion(client *frauddetector.Client, modelID, version string) {
	resp, err := client.DescribeModelVersions(context.TODO(), &frauddetector.DescribeModelVersionsInput{
		ModelId:            aws.String(modelID),
		ModelVersionNumber: aws.String(version),
		ModelType:          types.ModelTypeEnumTransactionFraudInsights,
	})
	if err != nil {
		log.Printf("Error describing model version: %v", err)
		return
	}

	for _, mv := range resp.ModelVersionDetails {
		fmt.Println("🔍 Model Version Details:")
		fmt.Printf("- Version: %s\n", aws.ToString(mv.ModelVersionNumber))
		fmt.Printf("- Status: %s\n", aws.ToString(mv.Status))
		fmt.Printf("- ARN: %s\n", aws.ToString(mv.Arn))
		fmt.Printf("- Training Data Source: %s\n", mv.TrainingDataSource)
		if mv.ExternalEventsDetail != nil {
			fmt.Printf("- Data Location: %s\n", aws.ToString(mv.ExternalEventsDetail.DataLocation))
		}
	}
}
