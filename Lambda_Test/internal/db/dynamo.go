package db

import (
	"Lambda_Test/internal/domain"
	"fmt"
)

// Not implemented
func InsertTransaction(transaction domain.Transaction) error {
	msg := fmt.Sprintf("Saved transaction with id %s and amount %.2f", transaction.Id, transaction.Amount)
	fmt.Println(msg)

	return nil
}
