package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/your-project/config"
	"github.com/your-project/models"
)

type App struct {
	config   *config.Config
	producer sarama.SyncProducer
	client   *http.Client
}

func NewApp(cfg *config.Config) (*App, error) {
	producer, err := createKafkaProducer(cfg.KafkaBrokers)
	if err != nil {
		return nil, fmt.Errorf("erreur création producteur Kafka: %w", err)
	}

	return &App{
		config:   cfg,
		producer: producer,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	log.Printf("🚀 Démarrage du scraper...")
	
	ticker := time.NewTicker(a.config.ScrapeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := a.scrapeAndSend(ctx); err != nil {
				log.Printf("❌ Erreur pendant le scraping: %v", err)
			}
		}
	}
}

func (a *App) scrapeAndSend(ctx context.Context) error {
	cryptos, err := a.fetchCryptoData(ctx)
	if err != nil {
		return fmt.Errorf("erreur récupération données: %w", err)
	}

	for _, crypto := range cryptos {
		if err := a.sendCryptoData(crypto); err != nil {
			log.Printf("❌ Erreur envoi données pour %s: %v", crypto.Symbol, err)
		}
	}
	return nil
}

func createKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	log.Printf("🔄 Configuration du producteur Kafka...")

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	log.Printf("📍 Tentative de connexion aux brokers: %v", brokers)

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("erreur création producteur après 5 tentatives: %v", err)
	}

	log.Printf("✅ Connexion à Kafka établie")
	return producer, nil
}

func (a *App) fetchCryptoData(ctx context.Context) ([]models.CryptoData, error) {
	resp, err := http.NewRequestWithContext(ctx, http.MethodGet, a.config.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err = a.client.Do(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.CoinCapResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (a *App) sendCryptoData(crypto models.CryptoData) error {
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
		Topic: a.config.KafkaTopic,
		Key:   sarama.StringEncoder(crypto.Symbol),
		Value: sarama.StringEncoder(jsonData),
	}

	partition, offset, err := a.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("erreur envoi message: %v", err)
	}

	log.Printf("✅ Données envoyées pour %s (prix: %.2f USD) sur partition %d à l'offset %d",
		crypto.Symbol, crypto.CurrentPrice, partition, offset)
	return nil
}

func main() {
	cfg := config.NewConfig()
	
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("❌ Erreur initialisation application: %v", err)
	}
	defer app.producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gestion gracieuse de l'arrêt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		cancel()
	}()

	if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("❌ Erreur fatale: %v", err)
	}
}
