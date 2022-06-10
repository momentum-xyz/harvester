package substrate

import (
	"reflect"
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	"github.com/stretchr/testify/assert"
)

func TestEraInfo(t *testing.T) {
	var cfg = harvester.Config{MQTT: mqtt.Config{Host: "localhost", Port: 1883, ClientId: "harvester"}}
	var chainCfg = harvester.ChainConfig{Name: "harvester", RPC: "ws://localhost:9944", ActiveTopics: []string{}}
	var h, _, _ = wire.NewHarvester(&cfg, func(err error) {})
	var sh, _ = NewHarvester(chainCfg, h.Publisher, h.Repository)

	t.Run("GetActiveEra", func(t *testing.T) {
		activeEra, err := sh.GetActiveEra()
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(activeEra).String(), reflect.Uint32.String())
	})

	t.Run("GetActiveEraDepth", func(t *testing.T) {
		activeEraDepth, err := sh.GetActiveEraDepth()
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(activeEraDepth), reflect.TypeOf([]byte("")))
	})

	t.Run("GetEraDepth", func(t *testing.T) {
		activeEra, _ := sh.GetActiveEra()
		eraDepth, err := sh.GetEraDepth(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(eraDepth), reflect.TypeOf([]byte("")))
	})

}
