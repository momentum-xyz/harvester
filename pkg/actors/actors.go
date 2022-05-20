package actors

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	exchangerateprovider "github.com/OdysseyMomentumExperience/harvester/pkg/exchange_rate_provider"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/OdysseyMomentumExperience/harvester/pkg/substrate"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	"github.com/oklog/run"
)

func Start(cfg harvester.Config, fn harvester.ErrorHandler) (*run.Group, error) {

	var err error

	h, cleanup, err := wire.NewHarvester(&cfg, fn)
	if err != nil {
		fn(err)
		cleanup()
		panic(err)

	}
	// wait for DB & MQTT to have connections
	time.Sleep(1 * time.Second)

	ctx := context.Background()
	g := new(run.Group)

	{
		ctx, cancel := context.WithCancel(ctx)
		g.Add(func() error {
			return exchangerateprovider.Start(ctx, fn, *h.Cfg, h.Publisher)
		}, func(err error) {
			fn(err)
			cleanup()
			cancel()
		})
	}

	var harvesters []harvester.ChainHarvester
	chains, err := getEnabledChains(cfg)

	if err != nil {
		return nil, err
	}

	for _, config := range chains {
		harvester, err := getChainHarvester(config, h.Publisher, h.Repository)
		if err != nil {
			fn(err)
			panic(err)
		}
		harvesters = append(harvesters, harvester)
	}

	for _, chainHarvester := range harvesters {

		{
			chainHarvester := chainHarvester
			ctx, cancel := context.WithCancel(ctx)
			g.Add(func() error {
				return chainHarvester.Start(ctx, fn)
			}, func(error) {
				fn(err)
				cleanup()
				cancel()
			})
		}
	}

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	return g, nil
}

func getEnabledChains(cfg harvester.Config) ([]harvester.ChainConfig, error) {
	var chains []harvester.ChainConfig
	for _, chainName := range cfg.EnabledChains {
		chainFound := false
		for _, chain := range cfg.Chains {
			if strings.EqualFold(chainName, chain.Name) {
				chains = append(chains, chain)
				chainFound = true
				break
			}
		}
		if !chainFound {
			return nil, fmt.Errorf("chain \"%s\" not supported", chainName)
		}
	}
	return chains, nil
}

func getChainHarvester(
	cfg harvester.ChainConfig,
	pub harvester.Publisher,
	repo harvester.Repository) (harvester.ChainHarvester, error) {
	switch cfg.Type {
	case "substrate":
		log.Infof("getting new substrate harvester for chain %s\n", cfg.Name)
		substratesvc, err := substrate.NewHarvester(cfg, pub, repo)
		if err != nil {
			return nil, err
		}
		return substratesvc, nil
	default:
		return nil, errors.New("could not recognize chain type")
	}
}
