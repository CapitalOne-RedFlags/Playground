package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func ImportEventsFromS3(client *frauddetector.Client, bucketName, fileName string) string {
	jobId := uuid.New().String()

	outputPath := fmt.Sprintf("s3://%s/output/%s", bucketName, jobId)
	iamRoleArn := "arn:aws:iam::920373029279:role/service-role/AmazonFraudDetector-DataAccessRole-1741553443122"

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
	return jobId
}

func WaitForBatchImportJobCompletion(client *frauddetector.Client, jobId string) {
	for {
		resp, err := client.GetBatchImportJobs(context.TODO(), &frauddetector.GetBatchImportJobsInput{
			JobId: aws.String(jobId),
		})
		if err != nil {
			log.Fatalf("Failed to get batch import job status: %v", err)
		}

		if len(resp.BatchImports) == 0 {
			log.Fatalf("No batch import jobs found with JobId %s", jobId)
		}

		status := resp.BatchImports[0].Status
		fmt.Println("Batch Import Job Status:", status)

		if status == "COMPLETE" {
			fmt.Println("Import job completed successfully.")
			break
		} else if status == "FAILED" {
			log.Fatalf("Import job failed: %v", aws.ToString(resp.BatchImports[0].FailureReason))

		}

		time.Sleep(10 * time.Second) // Wait before checking again
	}
}
