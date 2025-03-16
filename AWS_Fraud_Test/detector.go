package main

import (
	"context"
	"errors"
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
// or skips them if they already exist
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
			variableType: "IP_ADDRESS",
		},
		{
			name:         "transaction_amount",
			dataType:     types.DataTypeFloat,
			defaultValue: "0.0",
			description:  "Amount of the transaction",
			variableType: "NUMERIC",
		},
		{
			name:         "email_address",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Email address of the user",
			variableType: "EMAIL_ADDRESS",
		},
	}

	for _, v := range variables {
		// First check if the variable already exists
		variableName := aws.String(v.name)

		// Get the variable to see if it exists
		_, err := client.GetVariables(context.TODO(), &frauddetector.GetVariablesInput{
			Name: variableName,
		})

		if err == nil {
			// Variable exists, skip creation
			fmt.Printf("Variable %s already exists, skipping creation.\n", v.name)
			continue
		}

		// Check if error is something other than "variable not found"
		var resourceNotFoundException *types.ResourceNotFoundException
		if !errors.As(err, &resourceNotFoundException) {
			log.Printf("Error checking if variable %s exists: %v", v.name, err)
			continue
		}

		// Variable doesn't exist, create it
		_, err = client.CreateVariable(context.TODO(), &frauddetector.CreateVariableInput{
			Name:         variableName,
			DataType:     v.dataType,
			DefaultValue: aws.String(v.defaultValue),
			Description:  aws.String(v.description),
			VariableType: aws.String(v.variableType),
			DataSource:   types.DataSourceEvent,
		})

		if err != nil {
			log.Printf("Failed to create variable %s: %v", v.name, err)
			continue
		}

		fmt.Printf("Variable %s created successfully.\n", v.name)
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
