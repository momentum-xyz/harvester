package coingecko

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var cfg = harvester.Config{}
var h, _, _ = wire.NewHarvester(&cfg, func(err error) {})

func TestFetchExchangeRate(t *testing.T) {
	t.Run("fetch exchange rate", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"test_token": map[string]float64{"usd": float64(100)}})

		price, err := fetchExchangeRate(*h.Cfg, "test_token")
		assert.Nil(t, err)
		assert.IsType(t, price, &harvester.TokenPrice{})
		assert.Equal(t, price.TokenName, "test_token")
		assert.Equal(t, price.Currency, "usd")
	})

	t.Run("fetch exchange rate error", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(500).
			JSON(map[string]string{"message": "error"})

		price, err := fetchExchangeRate(*h.Cfg, "test_token")
		assert.Nil(t, price)
		assert.Error(t, err)
		assert.IsType(t, err, &json.UnmarshalTypeError{})
	})

	t.Run("fetch exchange rate request failure", func(t *testing.T) {
		defer gock.Off()

		gock.New("fake_url").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(100)}})

		price, err := fetchExchangeRate(*h.Cfg, "test_token")
		assert.Nil(t, price)
		assert.Error(t, err)
		assert.IsType(t, err, &url.Error{})
	})
}

func TestFetchKusamaTokenExchangeRate(t *testing.T) {
	t.Run("fetch kusama exchange rate", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(100)}})

		provider, _ := NewExchangeRateProvider(cfg)
		price, err := provider.FetchKusamaTokenExchangeRate()
		assert.Nil(t, err)
		assert.IsType(t, provider, &CoinGeckoExchangeRateProvider{})
		assert.IsType(t, price, &harvester.TokenPrice{})
		assert.Equal(t, price.TokenName, "kusama")
		assert.Equal(t, price.Currency, "usd")
	})
}
