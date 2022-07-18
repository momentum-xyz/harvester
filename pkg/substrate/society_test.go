package substrate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSocietyMembers(t *testing.T) {
	t.Run("ProcessSocietyMembers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.processSocietyMembers(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "society-members", 1*time.Second)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}
