package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Read variables from the environment
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqLogin := os.Getenv("RABBITMQ_LOGIN")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqUserQueue := os.Getenv("RABBITMQ_USER_QUEUE")

	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s/", rabbitmqLogin, rabbitmqPassword, rabbitmqHost)
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		rabbitmqUserQueue, // Queue name
		true,              // Durable
		false,             // Auto-delete
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		queue.Name, // Queue name
		"",         // Consumer tag (empty for auto-generated)
		true,       // Auto-acknowledge (message is automatically acknowledged once received)
		false,      // Exclusive (allow other consumers)
		false,      // No-local (don't get messages sent by this connection)
		false,      // No-wait (don't wait for server confirmation)
		nil,        // Arguments (none)
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	// Start receiving messages
	fmt.Println("Waiting for messages")
	for msg := range msgs {
		// Print the message body
		fmt.Printf("Received message: %s\n", msg.Body)
	}
}
