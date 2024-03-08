package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	c := client.NewClient(rpc.TestnetRPCEndpoint)
	sig, err := c.RequestAirdrop(
		context.TODO(),
		"3vZzcmwHTGf6cB5CwwhQuVqneM9nidMoTjhNQH4adhVP", // address
		1e9, // lamports (1 SOL = 10^9 lamports)
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
	}
	fmt.Println(sig)
}
