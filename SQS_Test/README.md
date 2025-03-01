# SQS Test Project

This Go project demonstrates how to **publish and consume** messages using **AWS SQS**. It includes:
- A **publisher** that sends transaction events to an SQS queue.
- A **consumer** that polls messages from the queue and processes them.
- A **utility package** to handle queue creation and message operations.

---

## 🛠 Prerequisites

1. **Install Go** (if not already installed)
   ```sh
   go version
   ```
   If Go is not installed, download it from [golang.org](https://go.dev/dl/).

2. **Install AWS CLI** (if not installed)
   ```sh
   aws --version
   ```
   If AWS CLI is not installed, download it from [AWS CLI Docs](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).

3. **Configure AWS Credentials**
   Run the following command and enter your AWS access credentials:
   ```sh
   aws configure
   ```
   You will be prompted for:
   - AWS Access Key ID
   - AWS Secret Access Key
   - Default region (e.g., `us-east-1`)
   - Output format (leave default as `json`)

---

## Project Structure

```
SQS_Test/
├── cmd/                 # Application entry points
│   ├── publisher/       # Sends messages to SQS
│   │   ├── main.go
│   ├── consumer/        # Reads messages from SQS
│   │   ├── main.go
├── config/              # Loads configuration (e.g., queue name)
│   ├── config.go
├── pkg/
│   ├── sqs/             # Contains SQS-related utilities
│   │   ├── client.go    # AWS SQS client setup
│   │   ├── producer.go  # Message publishing logic
│   │   ├── consumer.go  # Message consumption logic
│   │   ├── queue.go     # Queue creation & retrieval logic
├── go.mod               # Go module file
├── go.sum               # Dependency lock file
└── README.md            # Project documentation
```

---

## How to Run the Project

### **1️Set Up Environment Variables**
Create a **.env** file in the project root (optional, but recommended):
```
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
SQS_QUEUE_NAME=MyTestQueue
```
Or, manually set the environment variables in the terminal:

```sh
export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export SQS_QUEUE_NAME="MyTestQueue"
```

---

### **Install Dependencies**
Navigate to the project root and run:
```sh
go mod tidy
```

---

### **Start the Consumer**
In one terminal, run:
```sh
go run cmd/consumer/main.go
```
This process will **poll messages** from the SQS queue.

---

### **Start the Publisher**
In another terminal, run:
```sh
go run cmd/publisher/main.go
```
This will **send messages** to the queue.

---

### **Check Output**
- The **publisher** logs messages being sent.
- The **consumer** logs messages being processed.

---

## Cleanup (Delete the Queue)
If you want to **delete the SQS queue** after testing, use:
```sh
aws sqs delete-queue --queue-url $(aws sqs get-queue-url --queue-name <MyTestQueue> --query 'QueueUrl' --output text)
```

---

## Notes
- Ensure **AWS credentials are set correctly** before running the scripts.
- The queue will be **created dynamically** if it does not exist.
- Modify `config/config.go` to adjust settings as needed.

---


