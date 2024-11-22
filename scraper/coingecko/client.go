package coingecko

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "strings"
)

type CoinGeckoClient struct {
    baseURL    string
    httpClient *http.Client
}

type CryptoData struct {
    ID            string  `json:"id"`
    Symbol        string  `json:"symbol"`
    CurrentPrice  float64 `json:"current_price"`
    LastUpdated   string  `json:"last_updated"`
}

func NewCoinGeckoClient() *CoinGeckoClient {
    return &CoinGeckoClient{
        baseURL: "https://api.coingecko.com/api/v3",
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *CoinGeckoClient) GetCryptoPrices(cryptoIDs []string) ([]CryptoData, error) {
    url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&per_page=100&page=1&sparkline=false",
        c.baseURL, strings.Join(cryptoIDs, ","))
    
    fmt.Printf("Requête URL: %s\n", url)  // Log de l'URL

    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("erreur HTTP: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status code invalide: %d", resp.StatusCode)
    }

    var cryptoData []CryptoData
    if err := json.NewDecoder(resp.Body).Decode(&cryptoData); err != nil {
        return nil, fmt.Errorf("erreur de décodage JSON: %v", err)
    }

    if len(cryptoData) == 0 {
        return nil, fmt.Errorf("aucune donnée reçue")
    }

    return cryptoData, nil
} 