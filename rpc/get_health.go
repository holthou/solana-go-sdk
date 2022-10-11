package rpc

import "context"

// GetHealthResponse is a full raw rpc response of `getVersion`
type GetHealthResponse JsonRpcResponse[string]

// GetVersion returns the current solana versions running on the node
func (c *RpcClient) GetHealth(ctx context.Context) (JsonRpcResponse[string], error) {
	return c.processGetHealth(c.Call(ctx, "getHealth"))
}

func (c *RpcClient) processGetHealth(body []byte, rpcErr error) (res JsonRpcResponse[string], err error) {
	err = c.processRpcCall(body, rpcErr, &res)
	return
}
