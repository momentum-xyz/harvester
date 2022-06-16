package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/publisher"
	"github.com/stretchr/testify/assert"
)

func TestGetSlashes(t *testing.T) {
	err := mockSh.getSlashes(func(err error) {}, mockHarvester.PerformanceMonitorClient, "slashes")
	assert.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		ctx.Done()
		cancel()
	}()
	err = mockSh.processSlashes(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "slashes", 1*time.Second)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	t.Run("publishSlashEvent()", func(t *testing.T) {
		err := mockSh.publishSlashEvent(Slash{
			AccountID:    "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
			Amount:       100,
			EraIndex:     7,
			SessionIndex: 5,
			Type:         "validator",
		}, "slashes")
		assert.Nil(t, err)

		_mqttClient := mqtt.GetMQTTClient(&mqtt.Config{}, func(err error) {})
		_mockPublisher, _ := publisher.NewPublisher(_mqttClient)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, mockRepository, _mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)
		err = _mockSh.publishSlashEvent(Slash{
			AccountID:    "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
			Amount:       100,
			EraIndex:     7,
			SessionIndex: 5,
			Type:         "validator",
		}, "slashes")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to publish message")
	})
}
