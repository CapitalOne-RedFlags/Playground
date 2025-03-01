package main

import (
	"Lambda_Test/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.BatchProcessingHandler)
}
