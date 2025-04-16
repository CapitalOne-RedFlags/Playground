// package main

// import (
// 	"encoding/csv"
// 	"log"
// 	"os"
// )

// func main() {
// 	input := "labeled_transactions_clean.csv"
// 	output := "labeled_transactions_fixed.csv"

// 	inFile, err := os.Open(input)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer inFile.Close()

// 	reader := csv.NewReader(inFile)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Replace 'user' with 'customer' in ENTITY_TYPE column
// 	// Find index of ENTITY_TYPE
// 	header := records[0]
// 	var entityTypeIndex int = -1
// 	for i, col := range header {
// 		if col == "ENTITY_TYPE" {
// 			entityTypeIndex = i
// 			break
// 		}
// 	}
// 	if entityTypeIndex == -1 {
// 		log.Fatal("ENTITY_TYPE column not found")
// 	}

// 	for i := 1; i < len(records); i++ {
// 		if records[i][entityTypeIndex] == "user" {
// 			records[i][entityTypeIndex] = "customer"
// 		}
// 	}

// 	outFile, err := os.Create(output)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer outFile.Close()

// 	writer := csv.NewWriter(outFile)
// 	defer writer.Flush()

// 	err = writer.WriteAll(records)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("ENTITY_TYPE values updated. New file:", output)
// }

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	input := "labeled_transactions_clean.csv"
	output := "labeled_transactions.csv"

	inFile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	header := records[0]

	// Find index of columns we need
	var labelIndex int = -1
	var txDateIndex int = -1

	for i, col := range header {
		if col == "label" {
			labelIndex = i
		}
		if col == "transaction_date" {
			txDateIndex = i
		}
	}

	if labelIndex == -1 {
		log.Fatal("label column not found")
	}

	// Create new header with EXACTLY the headers AWS expects
	// Note: Order matters here - using the exact order from the error message
	newHeader := []string{
		"transaction_amount",
		"transaction_id",
		"transaction_date",
		"EVENT_LABEL",
		"ENTITY_ID",
		"account_balance",
		"email_address",
		"account_id",
		"transaction_duration",
		"ENTITY_TYPE",
		"transaction_type",
		"phone_number",
		"email",
		"LABEL_TIMESTAMP", // This was missing in previous script
		"EVENT_TIMESTAMP",
		"EVENT_ID",
		"ip_address",
		"location",
	}

	// Create new records array with transformed data
	newRecords := [][]string{newHeader}

	for i := 1; i < len(records); i++ {
		// Map label values (0,1) to AWS Fraud Detector labels (legit,fraud)
		label := "legit"
		if records[i][labelIndex] == "1" {
			label = "fraud"
		}

		// Create unique EVENT_ID
		eventId := fmt.Sprintf("event-%03d", i)

		// Create timestamp in ISO format
		timestamp := "2025-04-14T18:00:00Z" // Default

		// Use time from transaction_date if available
		if txDateIndex != -1 {
			txTime := records[i][txDateIndex]
			timestamp = fmt.Sprintf("2025-04-14T%s:00Z", convertTime(txTime))
		}

		// Create row matching EXACTLY the order of expected headers
		newRow := []string{
			getValue(records[i], header, "transaction_amount"), // transaction_amount
			getValue(records[i], header, "transaction_id"),     // transaction_id
			getValue(records[i], header, "transaction_date"),   // transaction_date
			label, // EVENT_LABEL
			"user-" + getValue(records[i], header, "account_id")[2:], // ENTITY_ID
			getValue(records[i], header, "account_balance"),          // account_balance
			getValue(records[i], header, "email"),                    // email_address
			getValue(records[i], header, "account_id"),               // account_id
			getValue(records[i], header, "transaction_duration"),     // transaction_duration
			"customer", // ENTITY_TYPE
			getValue(records[i], header, "transaction_type"), // transaction_type
			getValue(records[i], header, "phone_number"),     // phone_number
			getValue(records[i], header, "email"),            // email
			timestamp,                                        // LABEL_TIMESTAMP
			timestamp,                                        // EVENT_TIMESTAMP
			eventId,                                          // EVENT_ID
			getValue(records[i], header, "ip_address"),       // ip_address
			getValue(records[i], header, "location"),         // location
		}

		newRecords = append(newRecords, newRow)
	}

	outFile, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	err = writer.WriteAll(newRecords)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("CSV transformed for AWS Fraud Detector. New file:", output)
	log.Println("Headers match exactly what AWS Fraud Detector expects.")
}

// Helper function to safely get value from a record
func getValue(record []string, header []string, colName string) string {
	for i, col := range header {
		if col == colName {
			if i < len(record) {
				return record[i]
			}
		}
	}
	return "" // Column not found or index out of range
}

// Convert time format from "4:29:00 PM" to "16:29"
func convertTime(timeStr string) string {
	// Parse time string
	t, err := time.Parse("3:04:05 PM", timeStr)
	if err != nil {
		// If parsing fails, return a default
		return "18:00"
	}

	// Format as 24-hour time
	return t.Format("15:04")
}
