package client

import (
	"fmt"

	bittrex "github.com/toorop/go-bittrex"
)

type walletBtc struct {
	Price     float64
	Total     float64
	Available float64
	Estimated float64
}

func (w walletBtc) String() string {
	return fmt.Sprintf("Bitcoin: Price %0.2f$, Available: %0.8f, Total: %0.8f \n",
		w.Price, w.Available, w.Total)

}

func (w *walletBtc) EstimateBtc(altcoins Wallet) {
	balance := w.Total
	for _, coin := range altcoins {
		balance += coin.btcBalance()
	}
	w.Estimated = balance
}

func BtcUsdPrice() float64 {
	ticker, _ := client.GetTicker("USDT-BTC")
	return ticker.Last
}

func GetWalletBtc() walletBtc {
	priceCh := make(chan float64)
	balanceCh := make(chan bittrex.Balance)
	go func() {
		priceCh <- BtcUsdPrice()
	}()
	go func() {
		balance, _ := client.GetBalance("BTC")
		balanceCh <- balance
	}()
	price := <-priceCh
	balance := <-balanceCh
	return walletBtc{
		Price:     price,
		Total:     balance.Balance,
		Available: balance.Available,
	}
}
