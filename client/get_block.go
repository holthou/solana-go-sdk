package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedBlock for solana-core v1.6
// GetBlock returns identity and transaction information about a confirmed block in the ledger
func (c *Client) GetBlock(ctx context.Context, args ...interface{}) (*GetConfirmBlockResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmBlockResponse `json:"result"`
	}{}

	num := 0
	for {
		err := c.request(ctx, "getBlock", args, &res)
		if err != nil {
			return &GetConfirmBlockResponse{}, err
		}
		//块不可用时，重复获取
		if res.Error.Code == -32004 {
			if num >= 30 {
				logrus.Infof("getBlock break loop for time out! param:%v", args)
				break
			}
			num++
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if num > 0 {
		logrus.Infof("getBlock loop:%d param:%v", num, args)
	}

	//返回的res.Error不为空，正常逻辑，接口返回错误
	if !reflect.DeepEqual(res.Error, ErrorResponse{}) {
		return &res.Result, &res.Error
	}

	return &res.Result, nil
}
