package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
)

func main() {
	c := client.NewClient("https://fittest-summer-sheet.solana-mainnet.quiknode.pro/92548499ea403afb4c0c6c747d5a5727e7d3d2a4")
	nonceAccountAddr := "5jfGwxqrdtXkz6qYb78qHx1vsU3CjzHvQizjnFmXPF1D"
	nonceAccount, err := c.GetAccountInfo(context.Background(), nonceAccountAddr)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
	}

	fmt.Println(hex.EncodeToString(nonceAccount.Data))
	fmt.Printf("%+v\n", nonceAccount)
	/*
		type NonceAccount struct {
			Version          uint32
			State            uint32
			AuthorizedPubkey common.PublicKey
			Nonce            common.PublicKey
			FeeCalculator    FeeCalculator
		}
	*/
}
