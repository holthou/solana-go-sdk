package client

import (
	"context"
	"errors"
)

func (s *Client) GetSlot(ctx context.Context) (uint64, error) {
	arg := map[string]interface{}{
		"commitment": CommitmentFinalized,
	}

	res := struct {
		GeneralResponse
		Result uint64 `json:"result"`
	}{}
	err := s.request(ctx, "getSlot", []interface{}{arg}, &res)

	if err != nil {
		return 0, err
	}
	if res.Error != (ErrorResponse{}) {
		return 0, errors.New(res.Error.Message)
	}
	return res.Result, nil
}
