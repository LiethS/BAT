package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
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
	data []float64, sample bool
	
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}
	
	url := "https://localhost:3001/api/history"
	
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch historical data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	


	var history []HistoricalPrice
	for _, entry := range chart.Prices {
		timestamp := int64(entry[0]) / 1000
		date := time.Unix(timestamp, 0).Format("2006-01-02")
		history = append(history, HistoricalPrice{
			Date:  date,
			Price: entry[1],
		})
	}
	


	n := float64(len(data))
	if n == 0 {
		return 0
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
	
	if sample && n > 1 {
		variance /= (n - 1)
		} else {
			variance /= n
		}
		
		//final value of std
		return math.Sqrt(variance)
	}
<<<<<<< HEAD
	
	func zscore() {
		
	}
	
	func enableCORS(w *http.ResponseWriter) {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	// // ----- Manual valuation input -----
	// 	val := -1.87
	// 	fmt.Printf("Valuation metric: %.2f\n", val)
	
	// 	// Market condition evaluation
	// 	switch {
		// 	case val >= -2.0 && val < -1.5:
			// 		fmt.Println("Market Condition: EXTREME OVERSOLD")
			// 	case val >= -1.5 && val < -0.2:
				// 		fmt.Println("Market Condition: OVERSOLD")
				// 	case val >= -0.2 && val <= 0.2:
					// 		fmt.Println("Market Condition: NEUTRAL")
					// 	case val > 0.2 && val <= 1.5:
						// 		fmt.Println("Market Condition: OVERBOUGHT")
						// 	case val > 1.5 && val <= 2.0:
							// 		fmt.Println("Market Condition: EXTREME OVERBOUGHT")
							// 	default:
								// 		fmt.Println("Market Condition: Error - val out of expected range")
								// 	}
=======

	//final value of std
	return math.Sqrt(variance)
}

func zscore(value float64, data []float64, sample bool) float64 {
	n := float64(len(data))
	if n == 0 {
		return 0
	}

	//mean calculation
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / n
	std := StdDev(data, sample)

	if std == 0 {
		return 0
	}

	return (value - mean) / std

}
>>>>>>> 59b5d3cb349fe67cf85aad548cef7ea4feccef80
