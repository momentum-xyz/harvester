package harvester

type ExchangeRateProvider interface {
	FetchKusamaTokenExchangeRate() (*TokenPrice, error)
}

type TokenPrice struct {
	ProviderName string  `json:"providerName"`
	TokenName    string  `json:"tokenName"`
	Currency     string  `json:"currency"`
	Value        float64 `json:"value"`
}

type Topics struct {
	KusamaUSD   string `yaml:"kusamaUSD" envconfig:"KSM_USD"`
	PolkadotUSD string `yaml:"polkadotUSD" envconfig:"DOT_USD"`
}

// Config : structure to hold coingecko Provider configuration
type CoinGeckoConfig struct {
	Endpoint string `yaml:"endpoint" envconfig:"COINGECKO_ENDPOINT"`
	APIKey   string `yaml:"apiKey" envconfig:"COINGECKO_API_KEY"`
}

// Config : structure to hold Exchange Rate Provider configuration
type ExchangeRateProviderConfig struct {
	Active                      string          `yaml:"active" envconfig:"ACTIVE"`
	Topics                      Topics          `yaml:"topics"`
	PushInterval                int             `yaml:"pushInterval" envconfig:"PUSH_INTERVAL"`
	CoinGeko                    CoinGeckoConfig `yaml:"coingecko"`
	DisableExchangeRateProvider bool            `yaml:"disableExchangeRateProvider" envconfig:"DISABLE_EXCHANGE_RATE_PROVIDER"`
}
