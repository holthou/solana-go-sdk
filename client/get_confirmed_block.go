package client

import (
	"context"
	"reflect"
)

type GetConfirmBlockResponse struct {
	Blockhash         string                `json:"blockhash"`
	PreviousBlockhash string                `json:"previousBlockhash"`
	ParentSLot        uint64                `json:"parentSlot"`
	BlockTime         int64                 `json:"blockTime"`
	Transactions      []TransactionWithMeta `json:"transactions"`
	Rewards           []struct {
		Pubkey      string `json:"pubkey"`
		Lamports    int64  `json:"lamports"`
		PostBalance uint64 `json:"postBalance"`
		RewardType  string `json:"rewardType"` // type of reward: "fee", "rent", "voting", "staking"
	} `json:"rewards"`
}

func (s *Client) GetConfirmedBlock(ctx context.Context, args ...interface{}) (*GetConfirmBlockResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmBlockResponse `json:"result"`
	}{}
	err := s.request(ctx, "getConfirmedBlock", args, &res)
	if err != nil {
		return &GetConfirmBlockResponse{}, err
	}

	//返回的res.Error不为空，正常逻辑，接口返回错误
	if !reflect.DeepEqual(res.Error, ErrorResponse{}) {
		return &res.Result, &res.Error
	}

	return &res.Result, nil
}
