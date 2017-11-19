package client

import (
	"fmt"
	"os"
	"sync"

	"github.com/theherk/viper"
	bittrex "github.com/toorop/go-bittrex"
)

var once sync.Once
var client *bittrex.Bittrex

func GetClient() *bittrex.Bittrex {
	once.Do(func() {
		key := viper.GetString("key")
		secret := viper.GetString("secret")
		if key == "" || secret == "" {
			fmt.Println("Set API key and secret via config command.")
			os.Exit(1)
		}
		client = bittrex.New(key, secret)
	})
	return client
}
