package main

import (
	"fmt"
	"os"

	"github.com/sazor/bittrex-cli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// client := getBittrexClient()
	// wallet := walletBalances(client)
	// fmt.Println(wallet)
}
