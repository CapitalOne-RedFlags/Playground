package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
)

func main() {

	fmt.Println("Starting AWS Fraud Detector Setup...")

	// Get the AWS configuration
	cfg := GetAWSConfig()

	// Create a Fraud Detector client from the config
	client := frauddetector.NewFromConfig(cfg)

	// Step 1: Set up entity types and event types and labels
	CreateEntityType(client)
	CreateLabels(client)
	CreateEventVariables(client)
	CreateEventType(client)

	// Step 2: Create the fraud detector
	CreateDetector(client)

	// Step 3: Upload dataset to S3
	UploadToS3("redflags-bucket", "labeled_transactions.csv", "labeled_transactions.csv")

	jobId := ImportEventsFromS3(client, "redflags-bucket", "labeled_transactions.csv") // todos

	WaitForBatchImportJobCompletion(client, jobId)

	// Step 4: Create and train the Fraud Detector model
	CreateModel(client)
	TrainModel(client)

	// Step 5: Activate the trained model
	ActivateModel(client)

	// Step 6: Send an event to the fraud detector
	SendFraudEvent(client)

	// Step 7: Run a test prediction
	GetFraudScore(client)

	fmt.Println("AWS Fraud Detector Setup Completed.")
}
