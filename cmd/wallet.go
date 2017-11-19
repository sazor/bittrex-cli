// Copyright © 2017 Andrew Kozlov <andrewkozlov@icloud.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/sazor/bittrex-cli/client"
	"github.com/spf13/cobra"
)

var Order string
var OrderDirection string

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Display information about your wallet.",
	Long: `Print table with information about your altcoins:
			ticker, current price, average price of your position,
			price change in % and balance in btc.
			Also show display bitcoin balance and estimated balance of
			whole wallet in btc and $.`,
	Run: func(cmd *cobra.Command, args []string) {
		altWallet := client.WalletBalances(Order, OrderDirection)
		fmt.Println(altWallet)
		walletBtc := client.GetWalletBtc()
		walletBtc.EstimateBtc(altWallet)
		fmt.Println(walletBtc)
		fmt.Printf("Estimated: %0.8f => %0.2f$ \n", walletBtc.Estimated,
			walletBtc.Estimated*walletBtc.Price)
	},
}

func init() {
	RootCmd.AddCommand(walletCmd)
	walletCmd.Flags().StringVarP(&Order, "order", "o", "",
		"Choose table sorting (curprice, avgprice, change, balance)")
	walletCmd.Flags().StringVarP(&OrderDirection, "orderdir", "d", "asc",
		"Order direction (asc / desc)")
}
