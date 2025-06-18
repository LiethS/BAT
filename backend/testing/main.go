package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response from CoinGecko
type CoinGeckoResponse struct {
	MarketData struct {
		CurrentPrice map[string]float64 `json:"current_price"`
	} `json:"market_data"`
}

// Response to client
type PriceResponse struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

func main() {
	http.HandleFunc("/api/price", handlePrice)

	log.Println("✅ Server running at http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
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