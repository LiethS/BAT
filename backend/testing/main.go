package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type YahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamp []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func main() {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/BTC-USD?interval=1d&range=10y"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	// Spoof browser user-agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	var data YahooResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(data.Chart.Result) == 0 {
		log.Fatal("No results returned from Yahoo API")
	}

	timestamps := data.Chart.Result[0].Timestamp
	closes := data.Chart.Result[0].Indicators.Quote[0].Close

	if len(timestamps) == 0 || len(closes) == 0 {
		log.Fatal("No timestamps or close prices available")
	}

	latestIndex := len(timestamps) - 1
	latestTime := time.Unix(timestamps[latestIndex], 0).UTC()
	latestClose := closes[latestIndex]

	fmt.Printf("Latest BTC closing price on %s: $%.2f\n", latestTime.Format("2006-01-02"), latestClose)
}
