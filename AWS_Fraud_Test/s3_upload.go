package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// UploadToS3 uploads a file to an S3 bucket.
func UploadToS3(bucket, key, filePath string) {
	cfg := GetAWSConfig()

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("Failed to upload file to S3: %v", err)
	}
	fmt.Println("File uploaded to S3 successfully.")
}

func ImportEventsFromS3(client *frauddetector.Client, bucketName, fileName string) {
	// Generate a unique JobId
	jobId := uuid.New().String()

	// Define the OutputPath (this is where the result of the import job will be saved)
	outputPath := fmt.Sprintf("s3://%s/output/%s", bucketName, jobId)

	// Specify the IAM Role ARN (you need to create this role in your AWS account if it doesn't exist)
	iamRoleArn := "arn:aws:iam::123456789012:role/MyFraudDetectorRole" // Update with the correct IAM Role ARN

	// Set up import job for S3 CSV file
	_, err := client.CreateBatchImportJob(context.TODO(), &frauddetector.CreateBatchImportJobInput{
		EventTypeName: aws.String("transaction_event"),
		InputPath:     aws.String(fmt.Sprintf("s3://%s/%s", bucketName, fileName)),
		JobId:         aws.String(jobId),
		OutputPath:    aws.String(outputPath),
		IamRoleArn:    aws.String(iamRoleArn),
	})
	if err != nil {
		log.Fatalf("failed to create batch import job: %v", err)
	}
	fmt.Println("Batch import job created successfully.")
}
