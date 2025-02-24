package main
// This is my test line of code - Joe
import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Create a new SNS topic
func CreateSNSTopic(topicName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	input := &sns.CreateTopicInput{
		Name: aws.String(topicName),
	}

	result, err := snsClient.CreateTopic(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to create SNS topic: %v", err)
	}

	fmt.Println("SNS Topic Created:", *result.TopicArn)
	return *result.TopicArn, nil
}

// Subscribe an email or phone number to the topic
func SubscribeToSNSTopic(topicArn, protocol, endpoint string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	input := &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String(protocol),    // "email", "sms", "lambda", etc.
		Endpoint: aws.String(endpoint),    // email address or phone number
	}

	_, err = snsClient.Subscribe(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to subscribe to SNS topic: %v", err)
	}

	fmt.Printf("Successfully subscribed %s to SNS topic %s\n", endpoint, topicArn)
	return nil
}

// Publish a fraud alert message
func PublishFraudAlert(topicArn, message string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
	}

	_, err = snsClient.Publish(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send SNS alert: %v", err)
	}

	fmt.Println("Fraud alert sent successfully!")
	return nil
}

func main() {
	// Create SNS topic
	topicName := "FraudAlerts"
	topicArn, err := CreateSNSTopic(topicName)
	if err != nil {
		log.Fatalf("Error creating SNS topic: %v", err)
	}

	// Subscribe email
	err = SubscribeToSNSTopic(topicArn, "email", "leewenjie.wjl@gmail.com")
	if err != nil {
		log.Fatalf("Error subscribing email: %v", err)
	}

	// SMS number 
	err = SubscribeToSNSTopic(topicArn, "sms", "+1234567890")
	if err != nil {
		log.Fatalf("Error subscribing SMS: %v", err)
	}

	// Publish alert
	alertMessage := "Fraud detected on account #12345. Immediate action required!"
	err = PublishFraudAlert(topicArn, alertMessage)
	if err != nil {
		log.Fatalf("Error sending fraud alert: %v", err)
	}
}

