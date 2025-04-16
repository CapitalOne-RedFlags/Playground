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

// GetFraudScore evaluates a transaction for fraud.
func GetFraudScore(client *frauddetector.Client) {
	// Get the detector version
	versions, err := client.DescribeDetector(context.TODO(), &frauddetector.DescribeDetectorInput{
		DetectorId: aws.String("transaction_detector"),
	})
	if err != nil {
		log.Fatalf("Failed to get detector versions: %v", err)
	}
	if len(versions.DetectorVersionSummaries) == 0 {
		log.Fatalf("No detector versions found")
	}

	detectorVersionId := versions.DetectorVersionSummaries[0].DetectorVersionId
	fmt.Printf("Using detector version: %s\n", *detectorVersionId)

	// Create a suspicious transaction
	response, err := client.GetEventPrediction(context.TODO(), &frauddetector.GetEventPredictionInput{
		DetectorId:        aws.String("transaction_detector"),
		DetectorVersionId: detectorVersionId,
		EventId:           aws.String("test_event_001"),
		EventTypeName:     aws.String("transaction_event"),
		Entities: []types.Entity{
			{
				EntityType: aws.String("customer"),
				EntityId:   aws.String("test_customer_001"),
			},
		},
		EventTimestamp: aws.String(time.Now().UTC().Format("2006-01-02T15:04:05Z")),
		EventVariables: map[string]string{
			"transaction_id":       "TX001",
			"account_id":           "AC00123",
			"transaction_amount":   "2500.00", // High amount
			"transaction_date":     time.Now().Format("3:04:05 PM"),
			"transaction_type":     "Credit",
			"location":             "New York",
			"ip_address":           "192.168.1.1",
			"transaction_duration": "120",    // Long duration
			"account_balance":      "500.00", // Low balance compared to transaction
			"phone_number":         "1234567890",
			"email":                "test@example.com",
		},
	})
	if err != nil {
		log.Fatalf("Failed to get fraud score: %v", err)
	}

	// Display results
	fmt.Println("\nFraud Prediction Results:")
	fmt.Println("------------------------")
	for _, result := range response.ModelScores {
		fmt.Printf("Model: %s\n", *result.ModelVersion.ModelId)
		fmt.Printf("Fraud Score: %.2f\n", result.Scores["fraud_score"])
		fmt.Printf("Legit Score: %.2f\n", result.Scores["legit_score"])
	}

	// Display rule results if available
	if len(response.RuleResults) > 0 {
		fmt.Println("\nRule Results:")
		fmt.Println("-------------")
		for _, rule := range response.RuleResults {
			fmt.Printf("Rule: %s\n", *rule.RuleId)
			fmt.Printf("Outcome: %s\n", rule.Outcomes[0])
		}
	}
}
