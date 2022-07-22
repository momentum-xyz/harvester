package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestPendingExtrinsics(t *testing.T) {

	t.Run("ProcessPendingExtrinsics", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(5 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.ProcessPendingExtrinsics(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "extrinsics-pool")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("publishPendingExtrinsics", func(t *testing.T) {
		err := mockSh.publishPendingExtrinsics(func(err error) {}, mockHarvester.PerformanceMonitorClient, "extrinsics-pool")
		assert.Nil(t, err)
	})

	t.Run("ProcessExtrinsics", func(t *testing.T) {
		bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
		assert.Nil(t, err)
		amount := types.NewUCompactFromUInt(12345)
		c, err := types.NewCall(mockSh.metadata, "Balances.transfer", bob, amount)
		assert.Nil(t, err)
		ext := types.NewExtrinsic(c)
		extrinsics := []types.Extrinsic{ext}
		res := mockSh.ProcessExtrinsics(extrinsics)
		assert.Equal(t, len(res), 1)
		assert.IsType(t, res[0], harvester.Extrinsic{})
	})
}
