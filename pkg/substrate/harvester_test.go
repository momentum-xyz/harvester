package substrate

import (
	"reflect"
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var cfg = harvester.Config{
	MQTT: mqtt.Config{
		Host:     "localhost",
		Port:     1883,
		ClientId: "harvester",
	},
}
var chainCfg = harvester.ChainConfig{
	Name:         "harvester",
	RPC:          "ws://localhost:9944",
	ActiveTopics: []string{"block-creation-event"},
}

var h, _, _ = wire.NewHarvester(&cfg, func(err error) {})
var sh, _ = NewHarvester(chainCfg, h.Publisher, h.Repository)

func TestSubstrateHarvester(t *testing.T) {
	t.Run("NewHarvester()", func(t *testing.T) {
		_sh, err := NewHarvester(chainCfg, h.Publisher, h.Repository)
		assert.Nil(t, err)
		assert.IsType(t, _sh, &SubstrateHarvester{})

		_chainCfg := chainCfg
		_chainCfg.RPC = "ws://example.com"
		_sh, err = NewHarvester(_chainCfg, h.Publisher, h.Repository)
		assert.NotNil(t, err)
		assert.Nil(t, _sh)
	})

	t.Run("getApi()", func(t *testing.T) {
		api, err := getApi(chainCfg.RPC)
		assert.Nil(t, err)
		assert.IsType(t, api, &gsrpc.SubstrateAPI{})

		api, err = getApi("ws://example.com")
		assert.Nil(t, api)
		assert.NotNil(t, err)
	})

	t.Run("getMetadata()", func(t *testing.T) {
		api, err := getApi(chainCfg.RPC)
		assert.Nil(t, err)
		assert.IsType(t, api, &gsrpc.SubstrateAPI{})

		metadata, err := getMetadata(api)
		assert.IsType(t, metadata, &types.Metadata{})
		assert.Nil(t, err)
	})

	t.Run("getNetworkID()", func(t *testing.T) {
		assert.Equal(t, getNetworkID("kusama"), uint8(2))
		assert.Equal(t, getNetworkID("polkadot"), uint8(1))
		assert.Equal(t, getNetworkID("aaaa"), uint8(2))
	})

	t.Run("getActiveProcesses()", func(t *testing.T) {
		processes := sh.getActiveProcesses()
		assert.Equal(t, len(processes), 1)
		assert.Equal(t, reflect.ValueOf(processes[0]).Pointer(), reflect.ValueOf(sh.ProcessNewHeader).Pointer())
	})

	t.Run("topicProcessorStore()", func(t *testing.T) {
		process := sh.topicProcessorStore()("block-creation-event")
		assert.Equal(t, reflect.ValueOf(process).Pointer(), reflect.ValueOf(sh.ProcessNewHeader).Pointer())

		process = sh.topicProcessorStore()("aaaaa")
		assert.Nil(t, process)
	})

	t.Run("Stop()", func(t *testing.T) {
		sh.Stop()
	})

}
