package client

import "context"

// GetTransaction returns transaction details for a confirmed transaction
func (c *Client) GetTransaction(ctx context.Context, txhash string) (*GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := c.request(ctx, "getTransaction", []interface{}{txhash, "json"}, &res)
	if err != nil {
		return &GetConfirmedTransactionResponse{}, err
	}
	return &res.Result, nil
}
