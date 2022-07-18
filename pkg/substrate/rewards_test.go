package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErasRewardPoints(t *testing.T) {
	t.Run("processErasRewardPoints", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.processErasRewardPoints(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "reward-event", 1*time.Second)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("publishErasRewardPoints", func(t *testing.T) {
		err := mockSh.publishErasRewardPoints(func(err error) {}, mockHarvester.PerformanceMonitorClient, "reward-event")
		assert.Nil(t, err)
	})

	t.Run("GetEraReward", func(t *testing.T) {
		rewards, err := mockSh.GetEraReward(i32tob(0))
		assert.Nil(t, err)
		assert.IsType(t, rewards, EraRewardPoints{})
		rewards, err = mockSh.GetEraReward(i32tob(1111111111))
		assert.Nil(t, err)
		assert.Nil(t, rewards.Individual)
	})
}
