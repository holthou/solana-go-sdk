package client

import "context"

type GetTokenAccountsByOwnerResponse struct {
	Account GetAccountInfoResponse `json:"account"`
	Pubkey  string                 `json:"pubkey"`
}

type AccountDataByOwner struct {
	Parsed struct {
		Info struct {
			Delegate        string `json:"delegate"`
			DelegatedAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UiAmount       float64 `json:"uiAmount"`
				UiAmountString string  `json:"uiAmountString"`
			} `json:"delegatedAmount"`
			IsNative    bool   `json:"isNative"`
			Mint        string `json:"mint"`
			Owner       string `json:"owner"`
			State       string `json:"state"`
			TokenAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UiAmount       float64 `json:"uiAmount"`
				UiAmountString string  `json:"uiAmountString"`
			} `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

func (s *Client) GetTokenAccountsByOwner(ctx context.Context, base58Addr, mint string) (*[]GetTokenAccountsByOwnerResponse, error) {
	emptyResult := make([]GetTokenAccountsByOwnerResponse, 0)
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
		return &emptyResult, err
	}

	if len(res.Result.Value) > 0 {
		return &res.Result.Value, nil
	} else {
		return &emptyResult, nil
	}
}
