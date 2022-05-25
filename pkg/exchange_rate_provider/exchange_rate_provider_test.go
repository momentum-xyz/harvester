package exchangerateprovider

import (
	"context"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/coingecko"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var cfg = harvester.Config{
	MQTT: mqtt.Config{
		Host:     "localhost",
		Port:     1883,
		ClientId: "harvester",
	},
	ExchangeRateProvider: harvester.ExchangeRateProviderConfig{
		Active:       "coingecko",
		PushInterval: 2,
	},
}
var h, _, _ = wire.NewHarvester(&cfg, func(err error) {})

const ExchangeProviderNotSupported = "configured exchange rate provider not supported"

func TestGetExchangeRateProvider(t *testing.T) {
	t.Run("get exchange rate provider", func(t *testing.T) {
		provider, err := getExchangeRateProvider(cfg)
		assert.Nil(t, err)
		assert.IsType(t, provider, &coingecko.CoinGeckoExchangeRateProvider{})
	})

	t.Run("exchange rate provider not supported", func(t *testing.T) {
		provider, err := getExchangeRateProvider(harvester.Config{})
		assert.Nil(t, provider)
		assert.Equal(t, err.Error(), ExchangeProviderNotSupported)
	})
}

func TestPublishTokenPrice(t *testing.T) {
	t.Run("publish token price", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(500)}})

		provider, _ := getExchangeRateProvider(*h.Cfg)
		err := publishTokenPrice(provider, h.Publisher, "test topic")
		assert.Nil(t, err)
	})

	t.Run("unsupported protocol scheme", func(t *testing.T) {
		provider, _ := getExchangeRateProvider(cfg)
		err := publishTokenPrice(provider, nil, "test topic")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupported protocol scheme")
	})

	t.Run("exchange rate provider not supported", func(t *testing.T) {
		_cfg := cfg
		_cfg.ExchangeRateProvider.Active = "aaaa"
		provider, err := getExchangeRateProvider(_cfg)
		assert.Nil(t, provider)
		assert.Equal(t, err.Error(), ExchangeProviderNotSupported)
	})
}

func TestStart(t *testing.T) {
	t.Run("Start", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(500)}})

		h, _, _ := wire.NewHarvester(&cfg, func(err error) {})
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := Start(ctx, func(err error) {}, cfg, h.Publisher)
		assert.Equal(t, "context canceled", err.Error())
	})

	t.Run("Start: DisableExchangeRateProvider", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(500)}})

		h, _, _ := wire.NewHarvester(&cfg, func(err error) {})
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		_cfg := cfg
		_cfg.ExchangeRateProvider.DisableExchangeRateProvider = true
		err := Start(ctx, func(err error) {}, _cfg, h.Publisher)
		assert.Equal(t, "context canceled", err.Error())
	})

	t.Run("Start: Get exchange rate provider error", func(t *testing.T) {
		defer gock.Off()

		gock.New("(.*)").Get("/api/v3/simple/price").Reply(200).
			JSON(map[string]interface{}{"kusama": map[string]float64{"usd": float64(500)}})

		h, _, _ := wire.NewHarvester(&cfg, func(err error) {})
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		_cfg := cfg
		_cfg.ExchangeRateProvider.Active = "aaaa"
		err := Start(ctx, func(err error) {}, _cfg, h.Publisher)
		assert.Equal(t, err.Error(), ExchangeProviderNotSupported)
	})
}
