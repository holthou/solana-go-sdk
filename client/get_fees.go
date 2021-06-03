package client

import (
	"context"
)

type GetFeeResponse struct {
	Context struct {
		Slot int `json:"slot"`
	} `json:"context"`
	Value struct {
		Blockhash     string `json:"blockhash"`
		FeeCalculator struct {
			LamportsPerSignature int64 `json:"lamportsPerSignature"`
		} `json:"feeCalculator"`
		LastValidSlot int `json:"lastValidSlot"`
	} `json:"value"`
}

func (s *Client) GetFees(ctx context.Context) (int64, error) {
	res := struct {
		GeneralResponse
		Result GetFeeResponse `json:"result"`
	}{}

	err := s.request(ctx, "getFees", []interface{}{}, &res)
	if err != nil {
		return 0, err
	}

	return res.Result.Value.FeeCalculator.LamportsPerSignature, nil
}
