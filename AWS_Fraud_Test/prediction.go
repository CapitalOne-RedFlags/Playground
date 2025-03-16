package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
)

// GetFraudScore evaluates a transaction for fraud.
func GetFraudScore(client *frauddetector.Client) {
	response, err := client.GetEventPrediction(context.TODO(), &frauddetector.GetEventPredictionInput{
		DetectorId:        aws.String("transaction_detector"),
		DetectorVersionId: aws.String("1"), // Use your version ID
		EventId:           aws.String("event123"),
		EventTypeName:     aws.String("transaction_event"),
		Entities: []types.Entity{
			{
				EntityType: aws.String("customer"),
				EntityId:   aws.String("cust123"),
			},
		},
		EventTimestamp: aws.String("2023-03-16T12:00:00Z"),
		EventVariables: map[string]string{
			"ip_address":         "192.168.1.1",
			"transaction_amount": "500.00",
		},
	})
	if err != nil {
		log.Fatalf("Failed to get fraud score: %v", err)
	}

	// Display results
	fmt.Println("Fraud Prediction Results:")
	for _, result := range response.ModelScores {
		fmt.Printf("Model: %s, Score: %f\n", *result.ModelVersion.ModelId, result.Scores["fraud_score"])
	}

	// Display rule results if available
	for _, rule := range response.RuleResults {
		fmt.Printf("Rule: %s, Outcome: %s\n", *rule.RuleId, rule.Outcomes[0])
	}
}
