# Lambda Test Project

This Go project demonstrates how to consume **AWS SQS** messages using **AWS Lambda**

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

## How to Run the Project

### **1: Create AWS SQS Queue**
```
aws sqs create-queue --queue-name TransactionQueue
```
If you decide to name the queue something other than TransactionQueue,
set the queue name in enviroment variables:

```
export SQS_QUEUE_NAME="YourQueueName"
```

---

### **2: Create IAM Role to Allow Lambda to Read from SQS**
In **AWS IAM**, create a new role that includes permissions **AWSLambdaSQSQueueExecutionRole**

---

### **3: Create AWS Lambda Function**
In the AWS Lambda Dashboard, go to **Create a Function**
For Function name, put LambdaTest
For **Runtime**, choose **Amazon Linux 2023**
For **Architecture**, choose **x86_64**
For **Change default execution role**, choose **Use an existing role**, and select the IAM role you created in step 2
Click **Create function**

---

### **4: Add SQS Trigger**
Click on your function in the AWS Lambda Dashboard
If **SQS** is not yet listed as a trigger, click **Add Trigger** and select **SQS** from the dropdown
Select **TransactionQueue** as the **SQSQueue**
Use the default configuration for the trigger under **Event source mapping configuration**

---

### **5 Install Dependencies**
Navigate to the project root and run:
```sh
go mod tidy
```

---

### **6 Deploy Function Code to Lambda**
Run the following command to compile the code, zip the executable, and upload the zipfile to Lambda 

```sh
make deploy
```

---

### Running the Function
In **AWS SQS**, navigate to the queue you created and click **send or receive message**
Send a transaction of the form
```json
{"transaction_id": "1", "amount": 100.00}
```
In another tab, view your lambda function navigate to **Cloudwatch Logs**
The latest run's logs should show the transactions sent from SQS




