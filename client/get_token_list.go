package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//两个测试链，一个主链的id
const (
	MainnetBeta int = 101
	Testnet         = 102
	Devnet          = 103
)

type TokenInfo struct {
	ChainId    int           `json:"chainId"`
	Address    string        `json:"address"`
	Symbol     string        `json:"symbol"`
	Name       string        `json:"name"`
	Decimals   int           `json:"decimals"`
	LogoURI    string        `json:"logoURI"`
	Tags       []interface{} `json:"tags"`
	Extensions struct {
		Website      string `json:"website"`
		SerumV3Usdc  string `json:"serumV3Usdc"`
		SerumV3Usdt  string `json:"serumV3Usdt"`
		CoingeckoId  string `json:"coingeckoId"`
		Waterfallbot string `json:"waterfallbot"`
	} `json:"extensions"`
}

//获取指定链的Token list
func (s *Client) GetTokenList(ctx context.Context, envId int) ([]*TokenInfo, error) {
	result := struct {
		Name     string   `json:"name"`
		LogoURI  string   `json:"logoURI"`
		Keywords []string `json:"keywords"`
		Tags     struct {
		} `json:"tags"`
		Timestamp string      `json:"timestamp"`
		Tokens    []TokenInfo `json:"tokens"`
		Version   struct {
			Major int `json:"major"`
			Minor int `json:"minor"`
			Patch int `json:"patch"`
		} `json:"version"`
	}{}
	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", s.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	// http client and send request
	httpclient := &http.Client{}
	res, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
	}

	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return nil, fmt.Errorf("get status code: %d", res.StatusCode)
	}

	var list = make([]*TokenInfo, 0, len(result.Tokens))
	for _, t := range result.Tokens {
		v := t
		if t.ChainId == envId {
			list = append(list, &v)
		}
	}

	return list, nil
}
