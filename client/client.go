package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
)

type Client struct {
	RpcClient rpc.RpcClient
}

func New(opts ...rpc.Option) *Client {
	return &Client{
		RpcClient: rpc.New(opts...),
	}
}

func NewClient(endpoint string) *Client {
	return &Client{rpc.New(rpc.WithEndpoint(endpoint))}
}

type QuickSendTransactionParam struct {
	Instructions []types.Instruction
	Signers      []types.Account
	FeePayer     common.PublicKey
}

// Deprecated: please use sendTransaction
// QuickSendTransaction is a quick way to send tx
func (c *Client) QuickSendTransaction(ctx context.Context, param QuickSendTransactionParam) (string, error) {
	recentBlockhashRes, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash, err: %v", err)
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			Instructions:    param.Instructions,
			FeePayer:        param.FeePayer,
			RecentBlockhash: recentBlockhashRes.Blockhash,
		}),
		Signers: param.Signers,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new tx, err: %v", err)
	}
	rawTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize tx, err: %v", err)
	}
	res, err := c.RpcClient.SendTransactionWithConfig(
		ctx,
		base64.StdEncoding.EncodeToString(rawTx),
		rpc.SendTransactionConfig{Encoding: rpc.SendTransactionConfigEncodingBase64},
	)
	err = checkJsonRpcResponse(res, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

func checkJsonRpcResponse[T any](res rpc.JsonRpcResponse[T], err error) error {
	if err != nil {
		return err
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *Client) GetTokenAccountsByOwner(ctx context.Context, base58Addr string) (map[common.PublicKey]token.TokenAccount, error) {
	getTokenAccountsByOwnerResponse, err := c.RpcClient.GetTokenAccountsByOwnerWithConfig(
		ctx,
		base58Addr,
		rpc.GetTokenAccountsByOwnerConfigFilter{
			ProgramId: common.TokenProgramID.ToBase58(),
		},
		rpc.GetTokenAccountsByOwnerConfig{
			Encoding: rpc.AccountEncodingBase64,
		},
	)
	if err != nil {
		return nil, err
	}

	m := map[common.PublicKey]token.TokenAccount{}
	for _, v := range getTokenAccountsByOwnerResponse.Result.Value {
		accountInfo, err := c.rpcAccountInfoToClientAccountInfo(v.Account)
		if err != nil {
			return nil, err
		}
		tokenAccount, err := token.DeserializeTokenAccount(accountInfo.Data, accountInfo.Owner)
		if err != nil {
			return nil, err
		}
		m[common.PublicKeyFromString(v.Pubkey)] = tokenAccount
	}
	return m, err
}
