package client

import (
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env"
	bittrex "github.com/toorop/go-bittrex"
)

type config struct {
	BittrexKey    string `env:"BITTREX_KEY"`
	BittrexSecret string `env:"BITTREX_SECRET"`
}

var once sync.Once
var client *bittrex.Bittrex

func GetClient() *bittrex.Bittrex {
	once.Do(func() {
		cfg := config{}
		err := env.Parse(&cfg)
		if err != nil {
			fmt.Println("Set BITTREX_KEY and BITTREX_SECRET environment variables.")
			os.Exit(1)
		}
		client = bittrex.New(cfg.BittrexKey, cfg.BittrexSecret)
	})
	return client
}
