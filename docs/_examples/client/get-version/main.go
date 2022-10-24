package main

import (
	"context"
	"fmt"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
)

func main() {
	c := client.NewClient(rpc.MainnetRPCEndpoint)
	{
		resp, err := c.GetVersion(context.TODO())
		if err != nil {
			log.Fatalf("failed to version info, err: %v", err)
		}

		fmt.Println("version", resp.SolanaCore)
	}

	{
		arg := rpc.GetBlockConfig{
			Encoding:                       rpc.GetBlockConfigEncodingBase64, //对于这个接口,只可以用这个类型
			TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
			Rewards:                        pointer.Get[bool](false),
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: pointer.Get[uint8](0), //目前最新版本是0，可能会升级
		}

		resp, err := c.GetBlockWithConfig(context.TODO(), 155916607, arg)
		if err != nil {
			log.Fatalf("failed to version info, err: %v", err)
		}
		fmt.Println(resp)
	}

	{
		cfg := rpc.GetTransactionConfig{
			Encoding:                       rpc.TransactionEncodingBase64,
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: pointer.Get[uint8](0),
		}
		txid := "omVEuKY8h16SRFMZh1XxUxZNMPHdUqJoxTczzqFnWB3hQs1KPK5AAwtK91A5UAZRQ4Toh5N5ZkpZKoQFrYUwWF9"
		trx, err := c.GetTransactionWithConfig(context.TODO(), txid, cfg)
		if err != nil {
			log.Fatalf("failed to version info, err: %v", err)
		}

		fmt.Printf("%+v\n", *trx.Meta)

		if len(trx.Meta.LoadedAddresses.Writable) > 0 {
			for _, a := range trx.Meta.LoadedAddresses.Writable {
				trx.Transaction.Message.Accounts = append(trx.Transaction.Message.Accounts, common.PublicKeyFromString(a))
			}
		}

		if len(trx.Meta.LoadedAddresses.Readonly) > 0 {
			for _, a := range trx.Meta.LoadedAddresses.Readonly {
				trx.Transaction.Message.Accounts = append(trx.Transaction.Message.Accounts, common.PublicKeyFromString(a))
			}
		}
	}

}
