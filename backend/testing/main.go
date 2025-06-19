package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MarketChartResponse struct {
	Prices [][]float64 `json:"prices"`
}

type HistoricalPrice struct {
	Date  string  `json:"date"`  // e.g., "2025-06-18"
	Price float64 `json:"price"` // Close price
}

func main() {
	http.HandleFunc("/api/price", handlePrice)
	http.HandleFunc("/api/history", handleHistory)

	log.Println("✅ Server running at http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
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
		timestamp := int64(entry[0]) / 1000 // convert ms to seconds
		date := time.Unix(timestamp, 0).Format("2006-01-02")
		price := entry[1]

		history = append(history, HistoricalPrice{
			Date:  date,
			Price: price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func handlePrice(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	date := r.URL.Query().Get("date") // expected format: dd-mm-yyyy
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

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}