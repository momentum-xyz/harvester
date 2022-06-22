package substrate

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mysql"
	"github.com/OdysseyMomentumExperience/harvester/pkg/publisher"
	"github.com/OdysseyMomentumExperience/harvester/pkg/repository"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestBlocks(t *testing.T) {
	t.Run("ProcessNewHeader()", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.ProcessNewHeader(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "block-creation-event")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("ProcessFinalizedHeader()", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.ProcessFinalizedHeader(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "block-creation-event")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("fetchNewHead()", func(t *testing.T) {
		latestNewHead, err := mockSh.fetchNewHead(func(err error) {}, mockHarvester.PerformanceMonitorClient, 0, "block-creation-event")
		assert.Nil(t, err)
		assert.IsType(t, latestNewHead, types.NewU32(0))
	})

	t.Run("fetchNewHead(): fail to process header", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		_mqttClient := mqtt.GetMQTTClient(&mqtt.Config{}, func(err error) {})
		_mockPublisher, _ := publisher.NewPublisher(_mqttClient)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, mockRepository, _mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)

		latestNewHead, err := _mockSh.fetchNewHead(func(err error) {}, mockHarvester.PerformanceMonitorClient, 0, "block-creation-event")
		assert.Equal(t, latestNewHead, types.NewU32(0))
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error while processing new head")
	})

	t.Run("fetchFinalizedHead()", func(t *testing.T) {
		finalizedHead, err := mockSh.fetchFinalizedHead(func(err error) {}, mockHarvester.PerformanceMonitorClient, 0, "block-finalized-event")
		assert.Nil(t, err)
		assert.IsType(t, finalizedHead, types.NewU32(0))
	})

	t.Run("fetchFinalizedHead(): fail to process finalized head", func(t *testing.T) {
		finalizedHead, _ := mockSh.api.RPC.Chain.GetFinalizedHead()
		assert.IsType(t, finalizedHead, types.Hash{})

		_mqttClient := mqtt.GetMQTTClient(&mqtt.Config{}, func(err error) {})
		_mockPublisher, _ := publisher.NewPublisher(_mqttClient)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, mockRepository, _mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)

		latestNewHead, err := _mockSh.fetchFinalizedHead(func(err error) {}, mockHarvester.PerformanceMonitorClient, 0, "block-finalized-event")
		assert.Equal(t, latestNewHead, types.NewU32(0))
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error while processing finalized head")
	})

	t.Run("processHeader()", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		err := mockSh.processHeader(newHead, "block-creation-event", false)
		assert.Nil(t, err)
	})

	t.Run("processHeader(): publish header error", func(t *testing.T) {
		topic := "block-creation-event"
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		_mqttClient := mqtt.GetMQTTClient(&mqtt.Config{}, func(err error) {})
		_mockPublisher, _ := publisher.NewPublisher(_mqttClient)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, mockRepository, _mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)

		err := _mockSh.processHeader(newHead, topic, false)
		activeValidators, _ := mockSh.getCurrentSessionValidators()
		validatorID, _ := AccountIdToString(activeValidators[0])
		assert.Equal(t, reflect.TypeOf(activeValidators), reflect.SliceOf(reflect.TypeOf(types.NewAccountID([]byte("")))))

		errMsg := fmt.Errorf("publish header info for topic %v and validator %v is failed", topic, validatorID)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), errMsg.Error())
	})

	t.Run("publishHeader()", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		activeValidators, _ := mockSh.getCurrentSessionValidators()
		validatorID, _ := AccountIdToString(activeValidators[0])
		assert.Equal(t, reflect.TypeOf(activeValidators), reflect.SliceOf(reflect.TypeOf(types.NewAccountID([]byte("")))))

		err := mockSh.publishHeader(newHead, 1, validatorID, "block-creation-event", false)
		assert.Nil(t, err)
	})

	t.Run("publishHeader(): mqtt error", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		activeValidators, _ := mockSh.getCurrentSessionValidators()
		validatorID, _ := AccountIdToString(activeValidators[0])
		assert.Equal(t, reflect.TypeOf(activeValidators), reflect.SliceOf(reflect.TypeOf(types.NewAccountID([]byte("")))))

		_mqttClient := mqtt.GetMQTTClient(&mqtt.Config{}, func(err error) {})
		_mockPublisher, _ := publisher.NewPublisher(_mqttClient)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, mockRepository, _mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)

		err := _mockSh.publishHeader(newHead, 1, validatorID, "block-creation-event", false)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to publish message")
	})

	t.Run("publishHeader(): database error", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		assert.IsType(t, newHead, &types.Header{})

		activeValidators, _ := mockSh.getCurrentSessionValidators()
		validatorID, _ := AccountIdToString(activeValidators[0])
		assert.Equal(t, reflect.TypeOf(activeValidators), reflect.SliceOf(reflect.TypeOf(types.NewAccountID([]byte("")))))

		_mockDB, _, _ := mysql.NewDB(&mysql.Config{})
		_mockRepository, _ := repository.NewRepository(_mockDB, mockCfg.MySQL.Migrate)
		_mockHarvester, _ := harvester.NewHarvester(&mockCfg, _mockRepository, mockPublisher, mockPmc)
		_mockSh, _ := NewHarvester(mockChainCfg, _mockHarvester.Publisher, _mockHarvester.Repository)

		err := _mockSh.publishHeader(newHead, 1, validatorID, "block-creation-event", false)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "dial tcp")
	})

}
