package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"math"
)

// For /api/price endpoint
type CoinGeckoResponse struct {
	MarketData struct {
		CurrentPrice map[string]float64 `json:"current_price"`
	} `json:"market_data"`
}

type PriceResponse struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

// For /api/history endpoint
type MarketChartResponse struct {
	Prices [][]float64 `json:"prices"`
}

type HistoricalPrice struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

func main() {
	http.HandleFunc("/api/price", handlePrice)
	http.HandleFunc("/api/history", handleHistory)
	http.HandleFunc("/api/std", handleStd)
	http.HandleFunc("/api/zscore", handleZscore)

	log.Println("✅ Server running at http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func handlePrice(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Missing 'date' query param", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/bitcoin/history?date=%s", date)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch data from CoinGecko", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cg CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cg); err != nil {
		http.Error(w, "Error parsing CoinGecko response", http.StatusInternalServerError)
		return
	}

	price := cg.MarketData.CurrentPrice["usd"]
	result := PriceResponse{
		Date:  date,
		Price: price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	url := "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=30"

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch historical data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var chart MarketChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&chart); err != nil {
		http.Error(w, "Failed to parse CoinGecko response", http.StatusInternalServerError)
		return
	}

	var history []HistoricalPrice
	for _, entry := range chart.Prices {
		timestamp := int64(entry[0]) / 1000
		date := time.Unix(timestamp, 0).Format("2006-01-02")
		history = append(history, HistoricalPrice{
			Date:  date,
			Price: entry[1],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func handleStd(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	url := "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=30"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch historical data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var chart MarketChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&chart); err != nil {
		http.Error(w, "Failed to parse CoinGecko response", http.StatusInternalServerError)
		return
	}

	var data []float64
	for _, entry := range chart.Prices {
		// entry[0] is timestamp in ms, entry[1] is price
		data = append(data, entry[1])
	}

	n := float64(len(data))
	if n == 0 {
		return
	}

	//mean calculation
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / n

	//variance calculation
	var variance float64
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}

	sample := false
	if sample && n > 1 {
		variance /= (n - 1)
	} else {
		variance /= n
	}

	result := map[string]float64{
		"mean":   mean,
		"stddev": math.Sqrt(variance),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleZscore(w http.ResponseWriter, r *http.Request){
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	url := "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=30"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch historical data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var chart MarketChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&chart); err != nil {
		http.Error(w, "Failed to parse CoinGecko response", http.StatusInternalServerError)
		return
	}

	var data []float64
	for _, entry := range chart.Prices {
		// entry[0] is timestamp in ms, entry[1] is price
		data = append(data, entry[1])
	}

	n := float64(len(data))
	if n == 0 {
		http.Error(w, "No price data available", http.StatusInternalServerError)
		return
	}

	// mean
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / n

	// standard deviation (population)
	var variance float64
	for _, v := range data {
		variance += (v - mean) * (v - mean)
	}
	std := math.Sqrt(variance / n)
	if std == 0 {
		http.Error(w, "Standard deviation is zero", http.StatusInternalServerError)
		return
	}

	// Use the most recent price for z-score
	recentPrice := data[len(data)-1]
	z := (recentPrice - mean) / std

	// Response
	result := map[string]float64{
		"zscore": z,
		"price":  recentPrice,
		"mean":   mean,
		"stddev": std,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
