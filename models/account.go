package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Account represents a bank account.
type Account struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

// Transaction represents a transaction record.
type Transaction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	AccountID int                `json:"account_id" bson:"accountid"`
	Amount    float64            `json:"amount" bson:"amount"`
	Type      string             `json:"type" bson:"type"` // "deposit" or "withdrawal"
}

// TransactionRequest represents a transaction request from the API.
type TransactionRequest struct {
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
}
