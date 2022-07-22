package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	t.Run("Start", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.Start(ctx, mockHarvester.PerformanceMonitorClient, func(err error) {})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}
