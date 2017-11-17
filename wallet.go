package main

import (
	"fmt"
	"os"
	"sync"

	bittrex "github.com/toorop/go-bittrex"
)

const buyOrder = "LIMIT_BUY"

type WalletCoin struct {
	Balance   float64
	AvgPrice  float64
	CurrPrice float64
}

func (coin *WalletCoin) diffPrice() float64 {
	return coin.CurrPrice - coin.AvgPrice
}

func (coin *WalletCoin) percentDiffPrice() string {
	diff := (coin.diffPrice() / coin.AvgPrice) * 100
	return fmt.Sprintf("%0.2f %%", diff)
}

func (coin *WalletCoin) btcBalance() float64 {
	return coin.CurrPrice * coin.Balance
}
func calcAvgPrice(orders []bittrex.Order, posUnits float64) float64 {
	var totalCost, totalUnits float64
	for _, order := range orders {
		if order.OrderType == buyOrder {
			totalCost += order.Price + order.Commission
			totalUnits += order.Quantity - order.QuantityRemaining
			if totalUnits >= posUnits {
				break
			}
		}
	}
	return totalCost / totalUnits
}

func getNotNullTickers(client *bittrex.Bittrex, balances []bittrex.Balance) map[string]float64 {
	var tickers = make(map[string]float64)
	for _, coin := range balances {
		if coin.Balance > float64(0.0) && coin.Currency != "BTC" {
			tickers[coin.Currency] = coin.Balance
		}
	}
	return tickers
}

func walletBalances(client *bittrex.Bittrex) map[string]WalletCoin {
	balances, err := client.GetBalances()
	if err != nil {
		fmt.Println("Connection issues: %+v", err)
		os.Exit(1)
	}
	tickers := getNotNullTickers(client, balances)
	var wg sync.WaitGroup
	wg.Add(len(tickers))
	var wallet = make(map[string]WalletCoin)

	for ticker, balance := range tickers {
		go func(ticker string, balance float64) {
			tickerCh := make(chan bittrex.Ticker)
			go func() {
				ticker, _ := client.GetTicker("BTC-" + ticker)
				tickerCh <- ticker
			}()
			orderCh := make(chan []bittrex.Order)
			go func() {
				history, _ := client.GetOrderHistory("BTC-" + ticker)
				orderCh <- history
			}()
			currPrice, orderHistory := <-tickerCh, <-orderCh
			avgPrice := calcAvgPrice(orderHistory, balance)
			wallet[ticker] = WalletCoin{
				Balance:   balance,
				AvgPrice:  avgPrice,
				CurrPrice: currPrice.Last,
			}
			wg.Done()
		}(ticker, balance)
	}
	wg.Wait()
	return wallet
}
