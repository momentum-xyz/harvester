package harvester

import (
	"context"
	"time"

	influx_write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type PerformanceMonitorClient interface {
	// TODO add generics support
	WriteMetrics(*influx_write.Point, ErrorHandler)
	WriteProcessResponseMetrics(time.Time, string, ErrorHandler)
	WriteSystemRuntimeMetrics(context.Context, Config, ErrorHandler) error
}

//Config : structure to hold INFLUXDB configuration
type InfluxDbConfig struct {
	Url                   string `yaml:"url" envconfig:"INFLUXDB_URL"`
	Org                   string `yaml:"org" envconfig:"INFLUXDB_ORG"`
	Bucket                string `yaml:"bucket" envconfig:"INFLUXDB_BUCKET"`
	Token                 string `yaml:"token" envconfig:"INFLUXDB_TOKEN"`
	BatchSize             uint   `yaml:"batchSize" envconfig:"INFLUXDB_BATCH_SIZE"`
	RetryInterval         uint   `yaml:"retryInterval" envconfig:"INFLUXDB_RETRY_INTERVAL"`
	SystemMetricsInterval uint   `yaml:"systemMetricsInterval" envconfig:"SYSTEM_METRICS_INTERVAL"`
}
