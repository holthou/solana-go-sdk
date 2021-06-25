package client

import "context"

type GetTokenAccountsByOwnerResponse struct {
	Account GetAccountInfoResponse `json:"account"`
	Pubkey  string                 `json:"pubkey"`
}

func (s *Client) GetTokenAccountsByOwner(ctx context.Context, base58Addr, mint string) (*GetTokenAccountsByOwnerResponse, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context                           `json:"context"`
			Value   []GetTokenAccountsByOwnerResponse `json:"value"`
		} `json:"result"`
	}{}

	config := []interface{}{base58Addr,
		map[string]interface{}{"mint": mint},
		map[string]interface{}{"encoding": "jsonParsed"}}
	err := s.request(ctx, "getTokenAccountsByOwner", config, &res)
	if err != nil {
		return &GetTokenAccountsByOwnerResponse{}, err
	}

	//TODO 这里只获取数组中的第一个元素，是否合适
	if len(res.Result.Value) > 0 {
		return &res.Result.Value[0], nil
	} else {
		return &GetTokenAccountsByOwnerResponse{}, nil
	}
}
