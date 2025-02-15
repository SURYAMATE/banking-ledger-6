package queue

import (
	"encoding/json"
	"log"

	"banking-ledger-service/models" // Replace with your actual import path

	"github.com/streadway/amqp"
)

var rabbitMQConn *amqp.Connection
var rabbitMQChannel *amqp.Channel
var queueName = "transactions" // Define the queue name

// InitRabbitMQ initializes the RabbitMQ connection.
func InitRabbitMQ() {
	var err error
	rabbitMQConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // Modified connection string
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Error opening RabbitMQ channel: %v", err)
	}

	// Declare the queue
	_, err = rabbitMQChannel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring RabbitMQ queue: %v", err)
	}

	log.Println("Connected to RabbitMQ")
}

// PublishTransaction publishes a transaction request to the queue.
func PublishTransaction(transaction models.TransactionRequest) error {
	body, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = rabbitMQChannel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf("Published transaction: %v", transaction)
	return nil
}

func GetRabbitMQChannel() *amqp.Channel {
	return rabbitMQChannel
}

func GetQueueName() string {
	return queueName
}

// CloseRabbitMQ closes the RabbitMQ connection.  Important for cleanup.
func CloseRabbitMQ() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
}
