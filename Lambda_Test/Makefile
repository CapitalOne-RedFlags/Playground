# Executable must be named bootstrap using custom runtime in lambda
zipfile:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ./cmd/lambda/main.go
	zip batchProcessing.zip bootstrap

deploy: zipfile
	aws lambda update-function-code --function-name LambdaTest --zip-file fileb://batchProcessing.zip