package main

import (
	"fmt"

	"github.com/apcera/termtables"
)

func main() {
	client := getBittrexClient()
	wallet := walletBalances(client)

	table := termtables.CreateTable()
	table.AddHeaders("Ticker", "Current Price", "Avg Price", "Price Change", "BTC Balance")
	for ticker, walletCoin := range wallet {
		table.AddRow(ticker,
			fmt.Sprintf("%0.8f", walletCoin.CurrPrice),
			fmt.Sprintf("%0.8f", walletCoin.AvgPrice),
			walletCoin.percentDiffPrice(),
			fmt.Sprintf("%0.8f", walletCoin.btcBalance()))
	}
	fmt.Println(table.Render())
}
