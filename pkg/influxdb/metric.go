package influxdb

import (
	"context"
	"runtime"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influx_write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

func (ep *InfluxDb2Client) WriteMetrics(point *influx_write.Point, fn harvester.ErrorHandler) {
	ep.client.WritePoint(point)
}

func (ep *InfluxDb2Client) WriteProcessResponseMetrics(start time.Time, topic string, fn harvester.ErrorHandler) {
	p := influxdb2.NewPointWithMeasurement(topic).
		AddField("executionTime", time.Since(start).Milliseconds()).
		SetTime(time.Now())

	ep.client.WritePoint(p)
}

func (ep *InfluxDb2Client) WriteSystemRuntimeMetrics(ctx context.Context, cfg harvester.Config, fn harvester.ErrorHandler) error {
	pushInterval := cfg.InfluxDB.SystemMetricsInterval
	ticker := time.NewTicker(time.Second * time.Duration(pushInterval))
	defer ticker.Stop()
	var rtm runtime.MemStats

	for {

		select {
		case <-ticker.C:

			runtime.ReadMemStats(&rtm)

			lo := rtm.Mallocs - rtm.Frees
			p := influxdb2.NewPointWithMeasurement("system_metrics").
				AddField("go_goroutines", runtime.NumGoroutine()).
				AddField("go_mem_stats_alloc_bytes", rtm.Alloc).
				AddField("go_live_objects", lo).
				AddField("go_gc_duration", rtm.PauseTotalNs).
				AddField("go_gc_cycles", rtm.NumGC).
				SetTime(time.Now())

			ep.client.WritePoint(p)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
