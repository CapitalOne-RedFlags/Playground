package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
)

func main() {

	// // Open the original CSV
	// file, err := os.Open("bank_transactions_data.csv")
	// if err != nil {
	// 	log.Fatal("Error opening input file:", err)
	// }
	// defer file.Close()

	// reader := csv.NewReader(file)

	// // Read all rows
	// rows, err := reader.ReadAll()
	// if err != nil {
	// 	log.Fatal("Error reading CSV:", err)
	// }

	// if len(rows) == 0 {
	// 	log.Fatal("CSV is empty")
	// }

	// // Add 'Label' to the header
	// header := append(rows[0], "Label")

	// // Add label 0 to each data row
	// var updatedRows [][]string
	// updatedRows = append(updatedRows, header)
	// for _, row := range rows[1:] {
	// 	updatedRow := append(row, "0")
	// 	updatedRows = append(updatedRows, updatedRow)
	// }

	// // Write to a new file
	// outFile, err := os.Create("labeled_transactions.csv")
	// if err != nil {
	// 	log.Fatal("Error creating output file:", err)
	// }
	// defer outFile.Close()

	// writer := csv.NewWriter(outFile)
	// defer writer.Flush()

	// for _, row := range updatedRows {
	// 	if err := writer.Write(row); err != nil {
	// 		log.Fatal("Error writing row:", err)
	// 	}
	// }

	// log.Println("Successfully wrote labeled CSV!")

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
	UploadToS3("redflags-bucket", "labeled_transactions.csv", "bank_transactions_data.csv") // todos

	// ImportEventsFromS3(client, "redflags-bucket", "bank_transactions_data.csv") // todos

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
