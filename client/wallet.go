package client

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/apcera/termtables"
	bittrex "github.com/toorop/go-bittrex"
)

const buyOrder = "LIMIT_BUY"

type WalletCoin struct {
	Ticker    string
	Balance   float64
	AvgPrice  float64
	CurrPrice float64
}

type Wallet []WalletCoin

func (w Wallet) String() string {
	table := termtables.CreateTable()
	table.AddHeaders("Ticker", "Current Price", "Avg Price", "Price Change", "BTC Balance")
	for _, walletCoin := range w {
		table.AddRow(walletCoin.Ticker,
			fmt.Sprintf("%0.8f", walletCoin.CurrPrice),
			fmt.Sprintf("%0.8f", walletCoin.AvgPrice),
			fmt.Sprintf("%0.2f %%", walletCoin.percentDiffPrice()),
			fmt.Sprintf("%0.8f", walletCoin.btcBalance()))
	}
	return table.Render()
}

func (w Wallet) Sort(field string, direction string) {
	ascSort := direction != "desc"
	switch field {
	case "curprice":
		sortByCurrPrice(w, ascSort)
	case "avgprice":
		sortByAvgPrice(w, ascSort)
	case "change":
		sortByChange(w, ascSort)
	case "balance":
		sortByBalance(w, ascSort)
	default:
		sortByBalance(w, ascSort)
	}
}

func (coin *WalletCoin) diffPrice() float64 {
	return coin.CurrPrice - coin.AvgPrice
}

func (coin *WalletCoin) percentDiffPrice() float64 {
	return (coin.diffPrice() / coin.AvgPrice) * 100
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

func getNotNullTickers(balances []bittrex.Balance) map[string]float64 {
	var tickers = make(map[string]float64)
	for _, coin := range balances {
		if coin.Balance > float64(0.0) && coin.Currency != "BTC" {
			tickers[coin.Currency] = coin.Balance
		}
	}
	return tickers
}

func sortByCurrPrice(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].CurrPrice < wallet[j].CurrPrice
	})
}

func sortByAvgPrice(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].AvgPrice < wallet[j].AvgPrice
	})
}

func sortByBalance(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].btcBalance() < wallet[j].btcBalance()
	})
}

func sortByChange(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].percentDiffPrice() < wallet[j].percentDiffPrice()
	})
}

func fetchInfo(tickers map[string]float64) Wallet {
	var wg sync.WaitGroup
	wg.Add(len(tickers))

	var wallet Wallet
	for ticker, balance := range tickers {
		go func(ticker string, balance float64) {
			tickerCh := make(chan bittrex.Ticker)
			orderCh := make(chan []bittrex.Order)
			go func() {
				ticker, _ := client.GetTicker("BTC-" + ticker)
				tickerCh <- ticker
			}()
			go func() {
				history, _ := client.GetOrderHistory("BTC-" + ticker)
				orderCh <- history
			}()
			currPrice, orderHistory := <-tickerCh, <-orderCh
			avgPrice := calcAvgPrice(orderHistory, balance)
			wallet = append(wallet, WalletCoin{
				Balance:   balance,
				AvgPrice:  avgPrice,
				CurrPrice: currPrice.Last,
				Ticker:    ticker,
			})
			wg.Done()
		}(ticker, balance)
	}
	wg.Wait()
	return wallet
}

func WalletBalances(order string, orderdir string) Wallet {
	client := GetClient()
	balances, err := client.GetBalances()
	if err != nil {
		fmt.Println("Connection issues: %+v", err)
		os.Exit(1)
	}
	tickers := getNotNullTickers(balances)
	wallet := fetchInfo(tickers)
	wallet.Sort(order, orderdir)
	return wallet
}
