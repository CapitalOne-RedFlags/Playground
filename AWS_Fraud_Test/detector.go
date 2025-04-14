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
			name:        "1",
			description: "Fraudulent transaction",
		},
		{
			name:        "0",
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
		// List of your dataset fields as event variables
		{
			name:         "transaction_id",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Unique identifier for the transaction",
			variableType: "SESSION_ID", // Corrected to String
		},
		{
			name:         "account_id",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Unique identifier for the account",
			variableType: "USERAGENT", // Corrected to String
		},
		{
			name:         "transaction_amount",
			dataType:     types.DataTypeFloat,
			defaultValue: "0.0",
			description:  "Amount of the transaction",
			variableType: "PRICE", // Corrected to Numeric
		},
		{
			name:         "transaction_date",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Date when the transaction occurred",
			variableType: "FREE_FORM_TEXT", // Corrected to String
		},
		{
			name:         "transaction_type",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Type of the transaction (e.g., purchase, refund)",
			variableType: "PAYMENT_TYPE", // Corrected to String
		},
		{
			name:         "location",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Location of the transaction",
			variableType: "BILLING_CITY", // Corrected to String
		},
		// {
		// 	name:         "device_id",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "Device ID used for the transaction",
		// 	variableType: "String", // Corrected to String
		// },
		{
			name:         "ip_address",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "IP address of the user initiating the transaction",
			variableType: "BILLING_CITY", // Corrected to String
		},
		// {
		// 	name:         "merchant_id",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "ID of the merchant involved in the transaction",
		// 	variableType: "String", // Corrected to String
		// },
		// {
		// 	name:         "channel",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "Channel through which the transaction was made",
		// 	variableType: "String", // Corrected to String
		// },
		// {
		// 	name:         "customer_age",
		// 	dataType:     types.DataTypeInteger,
		// 	defaultValue: "0",
		// 	description:  "Age of the customer making the transaction",
		// 	variableType: "NUMERIC", // Corrected to Numeric
		// },
		// {
		// 	name:         "customer_occupation",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "Occupation of the customer",
		// 	variableType: "String", // Corrected to String
		// },
		{
			name:         "transaction_duration",
			dataType:     types.DataTypeFloat,
			defaultValue: "0.0",
			description:  "Duration of the transaction",
			variableType: "NUMERIC", // Corrected to Numeric
		},
		// {
		// 	name:         "login_attempts",
		// 	dataType:     types.DataTypeInteger,
		// 	defaultValue: "0",
		// 	description:  "Number of login attempts before the transaction",
		// 	variableType: "NUMERIC", // Corrected to Numeric
		// },
		{
			name:         "account_balance",
			dataType:     types.DataTypeFloat,
			defaultValue: "0.0",
			description:  "Balance of the customer's account at the time of transaction",
			variableType: "NUMERIC", // Corrected to Numeric
		},
		// {
		// 	name:         "previous_transaction_date",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "Date of the previous transaction",
		// 	variableType: "String", // Corrected to String
		// },
		{
			name:         "phone_number",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Phone number of the customer",
			variableType: "PHONE_NUMBER", // Corrected to String
		},
		{
			name:         "email",
			dataType:     types.DataTypeString,
			defaultValue: "",
			description:  "Email address of the customer",
			variableType: "EMAIL_ADDRESS", // Corrected to String
		},
		// {
		// 	name:         "transaction_status",
		// 	dataType:     types.DataTypeString,
		// 	defaultValue: "",
		// 	description:  "Status of the transaction (e.g., successful, failed)",
		// 	variableType: "String", // Corrected to String
		// },
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

func CreateEventType(client *frauddetector.Client) {
	_, err := client.PutEventType(context.TODO(), &frauddetector.PutEventTypeInput{
		Name: aws.String("transaction_event"),
		EventVariables: []string{
			"ip_address",         // keep this from original
			"transaction_amount", // keep this from original
			"email_address",
			"transaction_id", "account_id", "transaction_date", "transaction_type", "location", "transaction_duration", "account_balance", "phone_number", "email",

			// "TransactionID", "AccountID", "TransactionAmount", "TransactionDate",
			// "TransactionType", "Location", "DeviceID", "IP Address", "MerchantID",
			// "Channel", "CustomerAge", "CustomerOccupation", "TransactionDuration",
			// "LoginAttempts", "AccountBalance", "PreviousTransactionDate",
			// "Phone Number", "Email", "TransactionStatus",
		},
		EntityTypes: []string{"customer"},                 // Make sure this matches your defined entity
		Labels:      []string{"fraud", "legit", "0", "1"}, // Add this line to associate labels
		Tags:        []types.Tag{},
		// Optional: add Description, Inline, EventIngested, etc.
	})
	if err != nil {
		log.Fatalf("Failed to create event type: %v", err)
	}
	fmt.Println("Event type 'transaction_event' created successfully.")
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
