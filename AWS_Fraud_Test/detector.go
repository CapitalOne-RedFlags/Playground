package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
)

// CreateEntityType creates an entity type in AWS Fraud Detector.
func CreateEntityType(client *frauddetector.Client) {
	_, err := client.PutEntityType(context.TODO(), &frauddetector.PutEntityTypeInput{
		Name:        aws.String("customer"),
		Description: aws.String("Customer entity type"),
	})
	if err != nil {
		log.Fatalf("Failed to create entity type: %v", err)
	}
	fmt.Println("Entity type created.")
}

// CreateEventType creates an event type for fraud detection.
func CreateEventType(client *frauddetector.Client) {
	_, err := client.PutEventType(context.TODO(), &frauddetector.PutEventTypeInput{
		Name: aws.String("transaction_event"),
		EntityTypes: []string{ // ✅ Use string slice
			"customer",
		},
		EventVariables: []string{ // ✅ Use string slice
			"ip_address",
			"transaction_amount",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create event type: %v", err)
	}
	fmt.Println("Event type created.")
}

// CreateDetector creates a fraud detector.
func CreateDetector(client *frauddetector.Client) {
	_, err := client.PutDetector(context.TODO(), &frauddetector.PutDetectorInput{
		DetectorId:    aws.String("transaction_detector"),
		EventTypeName: aws.String("transaction_event"),
	})
	if err != nil {
		log.Fatalf("Failed to create detector: %v", err)
	}
	fmt.Println("Detector created.")
}
