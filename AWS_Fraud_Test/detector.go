package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
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

// CreateLabels creates the necessary labels before creating the event type
func CreateLabels(client *frauddetector.Client) {
	labels := []struct {
		name        string
		description string
	}{
		{
			name:        "fraud",
			description: "Fraudulent transaction",
		},
		{
			name:        "legit",
			description: "Legitimate transaction",
		},
	}

	for _, label := range labels {
		_, err := client.PutLabel(context.TODO(), &frauddetector.PutLabelInput{
			Name:        aws.String(label.name),
			Description: aws.String(label.description),
		})
		if err != nil {
			log.Printf("Failed to create label %s: %v", label.name, err)
			// Continue with the next label even if one fails
			continue
		}
		fmt.Printf("Label %s created successfully.\n", label.name)
	}
}

// CreateEventVariables creates the necessary event variables before creating the event type
func CreateEventVariables(client *frauddetector.Client) error {
	variables := []struct {
		name         string
		dataType     types.DataType
		defaultValue string
		description  string
		variableType string
	}{
		{
			name:         "ip_address",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "IP address of the user",
			variableType: "IP_ADDRESS", // Use valid variable type
		},
		{
			name:         "transaction_amount",
			dataType:     types.DataTypeFloat,
			defaultValue: "0.0",
			description:  "Amount of the transaction",
			variableType: "NUMERIC", // Use valid variable type for amounts
		},
		{
			name:         "email_address",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Email address of the user",
			variableType: "EMAIL_ADDRESS", // Use valid variable type
		},
	}

	successCount := 0
	for _, v := range variables {
		_, err := client.CreateVariable(context.TODO(), &frauddetector.CreateVariableInput{
			Name:         aws.String(v.name),
			DataType:     v.dataType,
			DefaultValue: aws.String(v.defaultValue),
			Description:  aws.String(v.description),
			VariableType: aws.String(v.variableType),
			DataSource:   types.DataSourceEvent,
		})
		if err != nil {
			log.Printf("Failed to create variable %s: %v", v.name, err)
			// Continue with the next variables even if one fails
			continue
		}
		fmt.Printf("Variable %s created successfully.\n", v.name)
		successCount++
	}

	if successCount < len(variables) {
		return fmt.Errorf("failed to create all variables, only %d of %d were successful",
			successCount, len(variables))
	}
	return nil
}

// CreateEventType creates an event type for fraud detection.
func CreateEventType(client *frauddetector.Client) {
	_, err := client.PutEventType(context.TODO(), &frauddetector.PutEventTypeInput{
		Name:           aws.String("transaction_event"),
		EntityTypes:    []string{"customer"},
		EventVariables: []string{"ip_address", "transaction_amount", "email_address"},
		Labels:         []string{"fraud", "legit"},
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
