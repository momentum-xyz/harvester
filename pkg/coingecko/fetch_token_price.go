package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

const (
	KSM      = "kusama"
	POLKADOT = "dot"
)

func (ep *CoinGeckoExchangeRateProvider) FetchKusamaTokenExchangeRate() (*harvester.TokenPrice, error) {

	return fetchExchangeRate(ep.cfg, KSM)
}

func fetchExchangeRate(cfg harvester.Config, tokenName string) (*harvester.TokenPrice, error) {
	baseURL := cfg.ExchangeRateProvider.CoinGeko.Endpoint

	url := fmt.Sprintf("%s/api/v3/simple/price", baseURL)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	q.Add("ids", tokenName)
	q.Add("vs_currencies", "usd")
	request.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tokens := make(map[string]TokenExchangeRate)

	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		return nil, err
	}
	return &harvester.TokenPrice{
		ProviderName: "CoinGecko",
		TokenName:    tokenName,
		Currency:     "usd",
		Value:        tokens[tokenName].USD,
	}, nil
}
