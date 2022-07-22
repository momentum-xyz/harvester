package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestChainSessionState(t *testing.T) {

	t.Run("ProcessChainSessionState", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.ProcessChainSessionState(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "session")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("currentSlotInSession", func(t *testing.T) {
		res, err := mockSh.currentSlotInSession(1)
		assert.Greater(t, res, types.NewU64(1))
		assert.Nil(t, err)
	})

	t.Run("currentSlotInEra", func(t *testing.T) {
		res, err := mockSh.currentSlotInEra(1)
		assert.GreaterOrEqual(t, res, types.NewU32(1))
		assert.Nil(t, err)
	})
}
