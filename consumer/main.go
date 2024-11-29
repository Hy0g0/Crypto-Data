package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
)

// Kafka config
var kafkaBrokers = []string{"kafka.kafka.svc.cluster.local:9092"}
var kafkaTopic = "matches"
var kafkaGroup = "matches"

// ClickHouse config
var ctx = context.Background()

// Define your message structure
type Message struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Details string `json:"details"`
}

func main() {
	// Set up Kafka consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		Topic:    kafkaTopic,
		GroupID:  kafkaGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer consumer.Close()

	fmt.Println("Consuming messages from Kafka...")

	// Set up ClickHouse connection using the DSN
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse-service.clickhouse.svc.cluster.local:9000"}, // Connection to the ClickHouse service
		Auth: clickhouse.Auth{
			Username: "default", // Use the correct username
			Password: "",        // Use the correct password
		},
		// The database is now part of the DSN, no longer a field in Options
	})
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	// Consume messages from Kafka and insert them into ClickHouse
	for {
		// Read message from Kafka
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		// Parse message into struct
		var message Message
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Insert data into ClickHouse
		err = insertIntoClickHouse(conn, message)
		if err != nil {
			log.Println("Error inserting into ClickHouse:", err)
		}
	}
}

// Function to insert message into ClickHouse
func insertIntoClickHouse(conn clickhouse.Conn, message Message) error {
	// Make sure the database is selected with the DSN

	query := fmt.Sprintf("INSERT INTO messages (id, name, details) VALUES (%d, '%s', '%s')",
		message.ID, message.Name, message.Details)

	// Execute query
	err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert into ClickHouse: %v", err)
	}

	fmt.Println("Inserted message into ClickHouse:", message)
	return nil
}
