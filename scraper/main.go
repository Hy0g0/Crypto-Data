package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
)

type CoinCapResponse struct {
	Data      []CryptoData `json:"data"`
	Timestamp int64        `json:"timestamp"`
}

type CryptoData struct {
	ID             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	CurrentPrice   float64 `json:"priceUsd,string"`
	MarketCap      float64 `json:"marketCapUsd,string"`
	Volume24h      float64 `json:"volumeUsd24Hr,string"`
	PriceChange24h float64 `json:"changePercent24Hr,string"`
	LastUpdated    string  `json:"timestamp"`
}

const (
	BASE_URL = "https://api.coincap.io/v2/assets"
	API_KEY  = "3282b43d-5f58-41de-a300-e049cc2fce0b" // Obtenez une clé sur https://coincap.io/api-key
)

func createKafkaProducer() (sarama.SyncProducer, error) {
	log.Printf("🔄 Configuration du producteur Kafka...")

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokers := []string{"kafka.kafka.svc.cluster.local:9092"}
	log.Printf("📍 Tentative de connexion aux brokers: %v", brokers)

	var producer sarama.SyncProducer
	var err error
	for i := 0; i < 5; i++ {
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			return producer, nil
		}
		log.Printf("Tentative %d échouée: %v. Nouvelle tentative dans 5 secondes...", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("erreur création producteur après 5 tentatives: %v", err)
}

func fetchCryptoData() ([]CryptoData, error) {
	resp, err := http.Get("https://api.coincap.io/v2/assets")
	if err != nil {
		time.Sleep(1 * time.Millisecond)
		return nil, err
	}
	defer resp.Body.Close()

	var response CoinCapResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		time.Sleep(1 * time.Millisecond)
		return nil, err
	}

	return response.Data, nil
}

func sendCryptoData(producer sarama.SyncProducer, crypto CryptoData) error {
	// Formater les données pour Kafka
	data := struct {
		Symbol         string  `json:"symbol"`
		Name           string  `json:"name"`
		Price          float64 `json:"price"`
		MarketCap      float64 `json:"marketCap"`
		Volume24h      float64 `json:"volume24h"`
		PriceChange24h float64 `json:"priceChange24h"`
		Timestamp      string  `json:"timestamp"`
	}{
		Symbol:         crypto.Symbol,
		Name:           crypto.Name,
		Price:          crypto.CurrentPrice,
		MarketCap:      crypto.MarketCap,
		Volume24h:      crypto.Volume24h,
		PriceChange24h: crypto.PriceChange24h,
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erreur marshalling JSON: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "crypto-prices",
		Key:   sarama.StringEncoder(crypto.Symbol),
		Value: sarama.StringEncoder(jsonData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("erreur envoi message: %v", err)
	}

	log.Printf("✅ Données envoyées pour %s (prix: %.2f USD) sur partition %d à l'offset %d",
		crypto.Symbol, crypto.CurrentPrice, partition, offset)
	return nil
}

func main() {
	log.Printf("🚀 Démarrage du scraper...")

	producer, err := createKafkaProducer()
	if err != nil {
		log.Fatalf("❌ Erreur initialisation Kafka: %v", err)
	}
	defer producer.Close()

	log.Printf("✅ Connexion à Kafka établie")

	// Boucle principale
	for {
		cryptos, err := fetchCryptoData()
		if err != nil {
			log.Printf("❌ Erreur récupération données: %v", err)
			time.Sleep(1 * time.Millisecond)
			continue
		}

		for _, crypto := range cryptos {
			if err := sendCryptoData(producer, crypto); err != nil {
				log.Printf("❌ Erreur envoi données pour %s: %v", crypto.Symbol, err)
			}
		}

		log.Printf("💤 Attente de 1ms avant la prochaine mise à jour...")
		time.Sleep(1 * time.Millisecond)
	}
}
