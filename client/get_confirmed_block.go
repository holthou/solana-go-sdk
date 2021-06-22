package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type GetConfirmBlockResponse struct {
	Blockhash         string                `json:"blockhash"`
	PreviousBlockhash string                `json:"previousBlockhash"`
	ParentSLot        uint64                `json:"parentSlot"`
	BlockTime         int64                 `json:"blockTime"`
	Transactions      []TransactionWithMeta `json:"transactions"`
	Signatures        []string              `json:"signatures"`
	Rewards           []struct {
		Pubkey      string `json:"pubkey"`
		Lamports    int64  `json:"lamports"`
		PostBalance uint64 `json:"postBalance"`
		RewardType  string `json:"rewardType"` // type of reward: "fee", "rent", "voting", "staking"
	} `json:"rewards"`
}

func (s *Client) getConfirmedBlockOld(ctx context.Context, args ...interface{}) (*GetConfirmBlockResponse, error) {
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

//封装借口，getConfirmedBlock解决 返回-32004的问题
func (s *Client) GetConfirmedBlock(ctx context.Context, args ...interface{}) (*GetConfirmBlockResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmBlockResponse `json:"result"`
	}{}

	num := 0
	for {
		err := s.request(ctx, "getConfirmedBlock", args, &res)
		if err != nil {
			return &GetConfirmBlockResponse{}, err
		}
		//块不可用时，重复获取
		if res.Error.Code == -32004 {
			if num >= 30 {
				logrus.Infof("GetConfirmedBlock break loop for time out! param:%v", args)
				break
			}
			num++
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if num > 0 {
		logrus.Infof("GetConfirmedBlock loop:%d param:%v", num, args)
	}

	//返回的res.Error不为空，正常逻辑，接口返回错误
	if !reflect.DeepEqual(res.Error, ErrorResponse{}) {
		return &res.Result, &res.Error
	}

	return &res.Result, nil
}
