package actors

import (
	"fmt"
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/stretchr/testify/assert"
)

var cfg = harvester.Config{
	Chains: []harvester.ChainConfig{{
		Name:      "kusama",
		RPC:       "wss://kusama-rpc.polkadot.io",
		Type:      "substrate",
		FromBlock: 0,
	}, {
		Name:      "polkadot",
		RPC:       "wss://rpc.polkadot.io",
		Type:      "substrate",
		FromBlock: 0,
	}},
	EnabledChains: []string{"kusama"},
}

func TestGetEnabledChains(t *testing.T) {
	t.Run("get enabled chains", func(t *testing.T) {
		chains, err := getEnabledChains(cfg)
		assert.Nil(t, err)
		assert.Equal(t, len(chains), 1)
		assert.Equal(t, chains[0].Name, cfg.Chains[0].Name)
		assert.IsType(t, chains, []harvester.ChainConfig{})

		cfg.EnabledChains = []string{"polkadot"}
		chains, err = getEnabledChains(cfg)
		assert.Nil(t, err)
		assert.Equal(t, len(chains), 1)
		assert.Equal(t, chains[0].Name, cfg.Chains[1].Name)
		assert.IsType(t, chains, []harvester.ChainConfig{})

		cfg.EnabledChains = []string{"kusama", "polkadot"}
		chains, err = getEnabledChains(cfg)
		assert.Nil(t, err)
		assert.Equal(t, len(chains), 2)
		assert.Equal(t, chains[0].Name, cfg.Chains[0].Name)
		assert.Equal(t, chains[1].Name, cfg.Chains[1].Name)
		assert.IsType(t, chains, []harvester.ChainConfig{})
	})

	t.Run("enabled chain not supported", func(t *testing.T) {
		_cfg := cfg
		_cfg.EnabledChains = []string{"aaa"}
		chains, err := getEnabledChains(_cfg)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("chain \"%s\" not supported", _cfg.EnabledChains[0]))
		assert.Nil(t, chains)
	})
}
