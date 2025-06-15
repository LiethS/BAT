package main

import (
	"fmt"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

const SYMBOL = "BINANCE:BTCUSDT" // https://www.tradingview.com/symbols/BTCUSDT/technicals/

func main() {
	var ta tradingview.TradingView
	
	// Fetch data for the specified symbol at a 4-hour interval
	err := ta.Get(SYMBOL, tradingview.Interval4Hour)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Get the summary trading recommendation
	recSummary := ta.Recommend.Global.Summary

	// Print the recommendation based on the signal
	switch recSummary {
	case tradingview.SignalStrongSell:
		fmt.Println("STRONG_SELL")
	case tradingview.SignalSell:
		fmt.Println("SELL")
	case tradingview.SignalNeutral:
		fmt.Println("NEUTRAL")
	case tradingview.SignalBuy:
		fmt.Println("BUY")
	case tradingview.SignalStrongBuy:
		fmt.Println("STRONG_BUY")
	default:
		fmt.Println("An error has occurred")
	}

	// Print the latest closing price
	clPrice := ta.Value.Prices.Close
	fmt.Println("Closing price:", clPrice)
}