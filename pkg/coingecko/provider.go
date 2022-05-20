package coingecko

import "github.com/OdysseyMomentumExperience/harvester/pkg/harvester"

type CoinGeckoExchangeRateProvider struct {
	cfg harvester.Config
}

func NewExchangeRateProvider(cfg harvester.Config) (*CoinGeckoExchangeRateProvider, error) {
	return &CoinGeckoExchangeRateProvider{
		cfg: cfg,
	}, nil
}
