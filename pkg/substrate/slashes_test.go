package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
	"github.com/stretchr/testify/assert"
)

func TestGetSlashes(t *testing.T) {
	err := sh.getSlashes(func(err error) {})
	assert.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		ctx.Done()
		cancel()
	}()
	err = sh.processSlashes(ctx, func(err error) {}, 1*time.Second)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	t.Run("publishSlashEvent()", func(t *testing.T) {
		err := sh.publishSlashEvent(Slash{
			AccountID:    "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
			Amount:       100,
			EraIndex:     7,
			SessionIndex: 5,
			Type:         "validator",
		})
		assert.Nil(t, err)

		_cfg := harvester.Config{MQTT: mqtt.Config{}}
		_h, _, _ := wire.NewHarvester(&_cfg, func(err error) {})
		_sh, _ := NewHarvester(chainCfg, _h.Publisher, _h.Repository)
		err = _sh.publishSlashEvent(Slash{
			AccountID:    "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
			Amount:       100,
			EraIndex:     7,
			SessionIndex: 5,
			Type:         "validator",
		})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to publish message")
	})
}
