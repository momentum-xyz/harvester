package performancemonitor

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/influxdb"
)

func NewPerformanceMonitorClient(cfg *harvester.Config, fn harvester.ErrorHandler) (harvester.PerformanceMonitorClient, error) {
	return influxdb.NewInfluxDBMonitorClient(cfg, fn)
}
