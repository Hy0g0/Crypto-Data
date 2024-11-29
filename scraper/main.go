package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/redis/go-redis/v9"
)

// Root Match Struct
// Match Struct
type Match struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	BeginAt          string      `json:"begin_at"`
	EndAt            *string     `json:"end_at,omitempty"`
	League           League      `json:"league"`
	Tournament       Tournament  `json:"tournament"`
	Serie            Serie       `json:"serie"`
	Videogame        Videogame   `json:"videogame"`
	Opponents        []Opponent  `json:"opponents"`
	Results          []Result    `json:"results"`
	StreamsList      []Stream    `json:"streams_list"`
	Games            []Game      `json:"games"`
	Status           string      `json:"status"`
	NumberOfGames    int         `json:"number_of_games"`
	Live             *Live       `json:"live,omitempty"`
	DetailedStats    bool        `json:"detailed_stats"`
	ModifiedAt       string      `json:"modified_at"`
	Forfeit          bool        `json:"forfeit"`
	Draw             bool        `json:"draw"`
	MatchType        string      `json:"match_type"`
	GameAdvantage    interface{} `json:"game_advantage,omitempty"`
	VideogameVersion *Version    `json:"videogame_version,omitempty"`
	Winner           *Winner     `json:"winner,omitempty"`
}

// League Struct
type League struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	URL        *string `json:"url,omitempty"` // made URL a pointer to handle null values
	ImageURL   string  `json:"image_url"`
	ModifiedAt string  `json:"modified_at"`
}

// Tournament Struct
type Tournament struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	BeginAt       string      `json:"begin_at"`
	EndAt         string      `json:"end_at"`
	LeagueID      int         `json:"league_id"`
	SerieID       int         `json:"serie_id"`
	PrizePool     interface{} `json:"prizepool,omitempty"`
	HasBracket    bool        `json:"has_bracket"`
	LiveSupported bool        `json:"live_supported"`
	Tier          string      `json:"tier"`
	ModifiedAt    string      `json:"modified_at"`
}

// Serie Struct
type Serie struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	BeginAt    string  `json:"begin_at"`
	EndAt      string  `json:"end_at"`
	WinnerID   *int    `json:"winner_id,omitempty"`
	WinnerType *string `json:"winner_type,omitempty"`
	ModifiedAt string  `json:"modified_at"`
	LeagueID   int     `json:"league_id"`
	Season     *string `json:"season,omitempty"`
	FullName   string  `json:"full_name"`
	Year       *int    `json:"year,omitempty"`
}

// Videogame Struct
type Videogame struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Opponent Struct
type Opponent struct {
	Type     string `json:"type"`
	Opponent Team   `json:"opponent"` // changed to Team type
}

// Team Struct (Previously Player)
type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Location   string `json:"location"`
	ModifiedAt string `json:"modified_at"`
	Acronym    string `json:"acronym"`
	ImageURL   string `json:"image_url"`
}

// Result Struct
type Result struct {
	TeamID int `json:"team_id"`
	Score  int `json:"score"`
}

// Stream Struct
type Stream struct {
	Main     bool   `json:"main"`
	Language string `json:"language"`
	EmbedURL string `json:"embed_url"`
	Official bool   `json:"official"`
	RawURL   string `json:"raw_url"`
}

// Game Struct
type Game struct {
	ID            int         `json:"id"`
	Position      int         `json:"position"`
	Status        string      `json:"status"`
	Length        interface{} `json:"length,omitempty"`
	BeginAt       *string     `json:"begin_at,omitempty"`
	EndAt         *string     `json:"end_at,omitempty"`
	Complete      bool        `json:"complete"`
	Finished      bool        `json:"finished"`
	DetailedStats bool        `json:"detailed_stats"`
	Forfeit       bool        `json:"forfeit"`
	MatchID       int         `json:"match_id"`
	WinnerType    *string     `json:"winner_type,omitempty"`
	Winner        *Winner     `json:"winner,omitempty"`
}

// Winner Struct
type Winner struct {
	ID   *int   `json:"id,omitempty"`
	Type string `json:"type"`
}

// Live Struct
type Live struct {
	Supported bool    `json:"supported"`
	URL       *string `json:"url,omitempty"`
	OpensAt   *string `json:"opens_at,omitempty"`
}

// Version Struct
type Version struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
}

const (
	PandaScoreAPIKey = "WphlC02gn0xLBkmxDeThNxkezOPPv9jF4U2fGSuQH0ZP-jIGGAg" // Replace with your API key
	PandaScoreURL    = "https://api.pandascore.co/matches"                   // Base URL
	KafkaTopic       = "matches"                                             // Kafka topic
	KafkaBroker      = "kafka.kafka.svc.cluster.local:9092"                  // Kafka broker
	RedisKey         = "pandascore:current_page"                             // Redis key
)

var ctx = context.Background()

func createKafkaProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{KafkaBroker}, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis-service.kafka.svc.cluster.local:6379", // Replace with your Redis address
		Password: "",                                           // No password set
		DB:       0,                                            // Use default DB
	})
}

func fetchMatches(page int) ([]Match, error) {
	url := fmt.Sprintf("%s?page=%d", PandaScoreURL, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", PandaScoreAPIKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch matches: %s", resp.Status)
	}

	var matches []Match
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func sendMatchToKafka(producer sarama.SyncProducer, match Match) error {
	jsonData, err := json.Marshal(match)
	if err != nil {
		return fmt.Errorf("failed to marshal match: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Value: sarama.StringEncoder(jsonData),
	}

	partition, offset, err := a.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send match to Kafka: %v", err)
	}

	log.Printf("✅ Match sent: %s (partition: %d, offset: %d)", match.Name, partition, offset)
	return nil
}

func getCurrentPage(redisClient *redis.Client) (int, error) {
	val, err := redisClient.Get(ctx, RedisKey).Result()
	if err == redis.Nil {
		return 1, nil // Default to page 1 if key doesn't exist
	} else if err != nil {
		return 0, err
	}

	page, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid page number in Redis: %v", err)
	}
	return page, nil
}

func setCurrentPage(redisClient *redis.Client, page int) error {
	return redisClient.Set(ctx, RedisKey, page, 0).Err()
}

func main() {
	log.Println("🚀 Starting Kafka Producer with Redis Tracking...")

	producer, err := createKafkaProducer()
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka producer: %v", err)
	}
	defer app.producer.Close()

	redisClient := createRedisClient()
	defer redisClient.Close()

	currentPage, err := getCurrentPage(redisClient)
	if err != nil {
		log.Fatalf("❌ Failed to get current page from Redis: %v", err)
	}

	log.Printf("🔄 Resuming from page %d...", currentPage)

	for {
		log.Printf("📡 Fetching matches from page %d...", currentPage)
		matches, err := fetchMatches(currentPage)
		if err != nil {
			log.Printf("❌ Error fetching matches: %v", err)
			break
		}

		if len(matches) == 0 {
			log.Println("🎉 All matches fetched and pushed to Kafka!")
			break
		}

		for _, match := range matches {
			if err := sendMatchToKafka(producer, match); err != nil {
				log.Printf("❌ Error sending match to Kafka: %v", err)
			}
		}

		currentPage++
		if err := setCurrentPage(redisClient, currentPage); err != nil {
			log.Printf("❌ Failed to update current page in Redis: %v", err)
			break
		}

		log.Printf("✅ Page %d processed. Moving to next page...", currentPage)
		time.Sleep(1 * time.Second) // Avoid rate-limiting by API
	}
}
