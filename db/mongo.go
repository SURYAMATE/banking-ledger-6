package db

import (
	"context"
	"log"

	"banking-ledger-service/models" // Replace with your actual import path

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var ledgerCollection *mongo.Collection

// InitMongo initializes the MongoDB connection.
func InitMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017") // Modified connection string
	var err error
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	ledgerCollection = mongoClient.Database("banking_ledger").Collection("transactions")
}

// LogTransaction logs a transaction to MongoDB.
func LogTransaction(transaction models.Transaction) error {
	_, err := ledgerCollection.InsertOne(context.Background(), transaction)
	return err
}

// GetTransactionHistory retrieves transaction history for an account from MongoDB.
func GetTransactionHistory(accountID int) ([]models.Transaction, error) {
	filter := bson.M{"accountid": accountID}
	cursor, err := ledgerCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []models.Transaction
	for cursor.Next(context.Background()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}
