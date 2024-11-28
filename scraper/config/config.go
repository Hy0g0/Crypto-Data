package config

type Config struct {
    BaseURL     string
    APIKey      string
    KafkaBrokers []string
    KafkaTopic   string
    ScrapeInterval time.Duration
}

func NewConfig() *Config {
    return &Config{
        BaseURL:        "https://api.coincap.io/v2/assets",
        APIKey:         "3282b43d-5f58-41de-a300-e049cc2fce0b",
        KafkaBrokers:   []string{"kafka.kafka.svc.cluster.local:9092"},
        KafkaTopic:     "crypto-prices",
        ScrapeInterval: time.Millisecond,
    }
} 