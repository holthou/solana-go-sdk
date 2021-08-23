package client

import "context"

//健康状态检查
func (s *Client) GetHealth(ctx context.Context) (string, error) {
	res := struct {
		GeneralResponse
		Result string `json:"result"`
	}{}
	err := s.request(ctx, "getHealth", []interface{}{}, &res)
	if err != nil {
		return "", err
	}
	if res.Result != "" && res.Result == "ok" {
		return res.Result, nil
	}
	return res.Error.Error(), nil
}
