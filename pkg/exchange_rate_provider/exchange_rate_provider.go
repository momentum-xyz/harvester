package exchangerateprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/coingecko"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
)

func Start(ctx context.Context,
	fn harvester.ErrorHandler,
	cfg harvester.Config,
	publisher harvester.Publisher,
	pmc harvester.PerformanceMonitorClient) error {

	if cfg.ExchangeRateProvider.DisableExchangeRateProvider {
		<-ctx.Done()
		return ctx.Err()
	}

	provider, err := getExchangeRateProvider(cfg)
	if err != nil {
		return err
	}

	pushInterval := cfg.ExchangeRateProvider.PushInterval * int(time.Second)

	log.Infof("scheduling kusama exchange rate provider for every %d sec", cfg.ExchangeRateProvider.PushInterval)

	ksmTopic := fmt.Sprintf("harvester/%s/ksm-usd-conversion-event", cfg.ExchangeRateProvider.Topics.KusamaUSD)

	ticker := time.NewTicker(time.Duration(pushInterval))

	for {
		select {
		case <-ctx.Done():
			log.Infof("Terminating exchange rate provider %v", ctx.Err())
			return ctx.Err()

		case <-ticker.C:
			err = publishTokenPrice(fn, provider, pmc, publisher, ksmTopic)
		}
		//TODO filter error types and handle individual errors
		if err != nil {
			fn(err)
		}

	}
}

func publishTokenPrice(fn harvester.ErrorHandler,
	provider harvester.ExchangeRateProvider,
	pmc harvester.PerformanceMonitorClient,
	publisher harvester.Publisher,
	ksmTopic string) error {

	defer pmc.WriteProcessResponseMetrics(time.Now(), ksmTopic, fn)
	var err error

	price, err := provider.FetchKusamaTokenExchangeRate()
	if err != nil {
		return err
	}
	priceJson, err := json.Marshal(price)
	if err != nil {
		return err
	}

	err = publisher.Publish(ksmTopic, string(priceJson))
	log.Debugf("publish kusama token usd price %f ", price.Value)

	if err != nil {
		return err
	}

	return err
}

func getExchangeRateProvider(cfg harvester.Config) (harvester.ExchangeRateProvider, error) {
	switch cfg.ExchangeRateProvider.Active {
	case "coingecko":
		log.Infof("getting coingecko exchange rate provider")

		provider, err := coingecko.NewExchangeRateProvider(cfg)
		if err != nil {
			return nil, err
		}
		return provider, nil
	default:
		return nil, errors.New("configured exchange rate provider not supported")
	}
}
