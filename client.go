package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	bittrex "github.com/toorop/go-bittrex"
)

type config struct {
	BittrexKey    string `env:"BITTREX_KEY"`
	BittrexSecret string `env:"BITTREX_SECRET"`
}

var client *bittrex.Bittrex

func getBittrexClient() *bittrex.Bittrex {
	if client != nil {
		return client
	}
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Set BITTREX_KEY and BITTREX_SECRET environment variables.")
		os.Exit(1)
	}
	client = bittrex.New(cfg.BittrexKey, cfg.BittrexSecret)
	return client
}
