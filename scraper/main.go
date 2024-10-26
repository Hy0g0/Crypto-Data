package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CryptoPrice struct {
	Crypto string `json:"crypto"`
	Price  string `json:"price"`
	Site   string `json:"site"`
}

func main() {
	// Kafka producer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafka.kafka.svc.cluster.local:9094",
		"client.id":         "crypto-price-producer",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}
	defer producer.Close()

	// Start the scraping and sending loop
	for {
		scrapeAndSend(producer)
		time.Sleep(1 * time.Second)
	}
}

func scrapeAndSend(producer *kafka.Producer) {
	url := "https://r.jina.ai/https://www.livecoinwatch.com/"

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// Define regex patterns for crypto name and price
	cryptoPattern := regexp.MustCompile(`\(?https:\/\/www\.livecoinwatch\.com\/price\/([A-Za-z]+)-.*\)`)
	pricePattern := regexp.MustCompile(`(?:https:\/\/www\.livecoinwatch\.com\/price\/[A-Za-z]+[-_].*\)\s*\|\s*\$)([0-9,.]+)`)

	// Find all matches
	cryptoMatches := cryptoPattern.FindAllStringSubmatch(string(body), -1)
	priceMatches := pricePattern.FindAllStringSubmatch(string(body), -1)

	// Store crypto prices in a map
	cryptoPrices := make(map[string]string)
	for i := range cryptoMatches {
		if len(cryptoMatches[i]) > 1 && len(priceMatches[i]) > 1 {
			cryptoName := cryptoMatches[i][1]
			price := priceMatches[i][1]
			cryptoPrices[cryptoName] = price
		}
	}

	// Send each crypto price to Kafka
	for crypto, price := range cryptoPrices {
		cryptoData := CryptoPrice{
			Crypto: crypto,
			Price:  price,
			Site:   "livecoin",
		}

		// Convert data to JSON
		message, err := json.Marshal(cryptoData)
		if err != nil {
			fmt.Printf("Failed to serialize data: %v\n", err)
			continue
		}

		// Produce message to Kafka
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &[]string{"crypto-prices"}[0], Partition: kafka.PartitionAny},
			Value:          message,
		}, nil)

		if err != nil {
			fmt.Printf("Failed to produce message: %v\n", err)
		} else {
			fmt.Printf("Sent message to crypto-prices topic: %s\n", message)
		}
	}

	// Ensure all messages are sent
	producer.Flush(15 * 1000)
}
