package main

import (
	"encoding/json"
	"fmt"
	"log"

	"banking-ledger-service/db"     // Replace with your actual import path
	"banking-ledger-service/models" // Replace with your actual import path
	"banking-ledger-service/queue"  // Replace with your actual import path

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Initialize database and queue connections
	db.InitPostgres()
	db.InitMongo()
	queue.InitRabbitMQ()
	defer queue.CloseRabbitMQ() // Close RabbitMQ connection when done

	// Get the channel and queue name from the queue package
	ch := queue.GetRabbitMQChannel()
	queueName := queue.GetQueueName()

	// Consume messages from the queue
	messages, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Error consuming from queue: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var transactionRequest models.TransactionRequest
			err := json.Unmarshal(d.Body, &transactionRequest)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue // Skip invalid messages
			}

			// Process the transaction
			err = ProcessTransaction(transactionRequest)
			if err != nil {
				log.Printf("Error processing transaction: %v", err)
				// Consider adding a retry mechanism or dead-letter queue for failed transactions
			}
		}
	}()

	log.Println("Transaction processor worker started.  Listening for messages...")
	<-forever // Keep the worker running
}

// ProcessTransaction processes a transaction, updating the account balance and logging the transaction.
func ProcessTransaction(transactionRequest models.TransactionRequest) error {
	log.Printf("Processing transaction: %v", transactionRequest)

	accountID := transactionRequest.AccountID
	amount := transactionRequest.Amount
	transactionType := transactionRequest.Type

	// 1. Get current account balance from PostgreSQL
	currentBalance, err := db.GetAccountBalance(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account balance: %w", err)
	}

	// 2. Calculate new balance
	newBalance := currentBalance
	if transactionType == "deposit" {
		newBalance += amount
	} else if transactionType == "withdrawal" {
		if currentBalance < amount {
			return fmt.Errorf("insufficient funds for withdrawal")
		}
		newBalance -= amount
	} else {
		return fmt.Errorf("invalid transaction type: %s", transactionType) // Should be caught by API validation
	}

	// 3. Update account balance in PostgreSQL
	err = db.UpdateAccountBalance(accountID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// 4. Log transaction in MongoDB
	transactionLog := models.Transaction{
		ID:        primitive.NewObjectID(),
		AccountID: accountID,
		Amount:    amount,
		Type:      transactionType,
	}
	err = db.LogTransaction(transactionLog)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}

	log.Printf("Transaction processed successfully. Account %d new balance: %f", accountID, newBalance)
	return nil
}
