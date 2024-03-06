package client

import (
	"context"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
)

type TokenAccount struct {
	token.TokenAccount
	PublicKey common.PublicKey
}

// TOOD 待测试函数的功能
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
		accountInfo, err := convertAccountInfo(v.Account)
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

func (c *Client) GetTokenAccountsByOwnerByMint(ctx context.Context, owner, mintAddr string) ([]TokenAccount, error) {
	return process(
		func() (rpc.JsonRpcResponse[rpc.ValueWithContext[rpc.GetProgramAccounts]], error) {
			return c.RpcClient.GetTokenAccountsByOwnerWithConfig(
				ctx,
				owner,
				rpc.GetTokenAccountsByOwnerConfigFilter{
					Mint: mintAddr,
				},
				rpc.GetTokenAccountsByOwnerConfig{
					Encoding: rpc.AccountEncodingBase64,
				},
			)
		},
		convertGetTokenAccountsByOwner,
	)
}

func (c *Client) GetTokenAccountsByOwnerByProgram(ctx context.Context, owner, programId string) ([]TokenAccount, error) {
	return process(
		func() (rpc.JsonRpcResponse[rpc.ValueWithContext[rpc.GetProgramAccounts]], error) {
			return c.RpcClient.GetTokenAccountsByOwnerWithConfig(
				ctx,
				owner,
				rpc.GetTokenAccountsByOwnerConfigFilter{
					ProgramId: programId,
				},
				rpc.GetTokenAccountsByOwnerConfig{
					Encoding: rpc.AccountEncodingBase64,
				},
			)
		},
		convertGetTokenAccountsByOwner,
	)
}

func (c *Client) GetTokenAccountsByOwnerWithContextByMint(ctx context.Context, owner, mintAddr string) (rpc.ValueWithContext[[]TokenAccount], error) {
	return process(
		func() (rpc.JsonRpcResponse[rpc.ValueWithContext[rpc.GetProgramAccounts]], error) {
			return c.RpcClient.GetTokenAccountsByOwnerWithConfig(
				ctx,
				owner,
				rpc.GetTokenAccountsByOwnerConfigFilter{
					Mint: mintAddr,
				},
				rpc.GetTokenAccountsByOwnerConfig{
					Encoding: rpc.AccountEncodingBase64,
				},
			)
		},
		convertGetTokenAccountsByOwnerAndContext,
	)
}

func (c *Client) GetTokenAccountsByOwnerWithContextByProgram(ctx context.Context, owner, programId string) (rpc.ValueWithContext[[]TokenAccount], error) {
	return process(
		func() (rpc.JsonRpcResponse[rpc.ValueWithContext[rpc.GetProgramAccounts]], error) {
			return c.RpcClient.GetTokenAccountsByOwnerWithConfig(
				ctx,
				owner,
				rpc.GetTokenAccountsByOwnerConfigFilter{
					ProgramId: programId,
				},
				rpc.GetTokenAccountsByOwnerConfig{
					Encoding: rpc.AccountEncodingBase64,
				},
			)
		},
		convertGetTokenAccountsByOwnerAndContext,
	)
}

func convertGetTokenAccountsByOwner(v rpc.ValueWithContext[rpc.GetProgramAccounts]) ([]TokenAccount, error) {
	tokenAccounts := make([]TokenAccount, 0, len(v.Value))
	for _, v := range v.Value {
		accountInfo, err := convertAccountInfo(v.Account)
		if err != nil {
			return nil, err
		}
		tokenAccount, err := token.DeserializeTokenAccount(accountInfo.Data, accountInfo.Owner)
		if err != nil {
			return nil, err
		}
		tokenAccounts = append(tokenAccounts, TokenAccount{
			TokenAccount: tokenAccount,
			PublicKey:    common.PublicKeyFromString(v.Pubkey),
		})
	}
	return tokenAccounts, nil
}

func convertGetTokenAccountsByOwnerAndContext(v rpc.ValueWithContext[rpc.GetProgramAccounts]) (rpc.ValueWithContext[[]TokenAccount], error) {
	tokenAccounts, err := convertGetTokenAccountsByOwner(v)
	if err != nil {
		return rpc.ValueWithContext[[]TokenAccount]{}, err
	}
	return rpc.ValueWithContext[[]TokenAccount]{
		Context: v.Context,
		Value:   tokenAccounts,
	}, nil
}
