package substrate

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestChainEraState(t *testing.T) {

	t.Run("processChainEraState", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(2 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.processChainEraState(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "era", 1*time.Second)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("publishEraInfo", func(t *testing.T) {
		err := mockSh.publishEraInfo(func(err error) {}, mockHarvester.PerformanceMonitorClient, "era")
		assert.Nil(t, err)
	})

	t.Run("GetErasValidatorReward", func(t *testing.T) {
		activeEra, _ := mockSh.GetActiveEra()
		erasValidatorReward, err := mockSh.GetErasValidatorReward(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(erasValidatorReward), reflect.TypeOf(types.NewU128(*big.NewInt(0))))
	})

	t.Run("getEraTotalStake", func(t *testing.T) {
		activeEra, _ := mockSh.GetActiveEra()
		totalStakeInEra, err := mockSh.getEraTotalStake(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(totalStakeInEra), reflect.TypeOf(types.NewU128(*big.NewInt(0))))
	})

	t.Run("getStakingRatio", func(t *testing.T) {
		activeEra, _ := mockSh.GetActiveEra()
		stakingRatio, err := mockSh.getStakingRatio(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(stakingRatio), reflect.TypeOf(float32(0)))
	})
}
