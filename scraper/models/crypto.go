package models

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