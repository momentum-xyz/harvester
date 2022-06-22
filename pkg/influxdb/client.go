package influxdb

import (
	"context"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	influx_db2 "github.com/influxdata/influxdb-client-go/v2"
	influx_api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
)

type InfluxDb2Client struct {
	client influx_api.WriteAPI
}

func NewInfluxDBMonitorClient(cfg *harvester.Config, fn harvester.ErrorHandler) (*InfluxDb2Client, error) {

	influxdb := cfg.InfluxDB

	influxClient := influx_db2.NewClientWithOptions(influxdb.Url, influxdb.Token,
		influx_db2.DefaultOptions().
			SetBatchSize(influxdb.BatchSize).
			SetPrecision(time.Millisecond).
			SetRetryInterval(influxdb.RetryInterval))

	_, err := influxClient.Health(context.Background())

	if err != nil {
		return nil, errors.Wrap(err, "error occurred while creating influxdb client")

	}

	nonBlockingClient := influxClient.WriteAPI(influxdb.Org, influxdb.Bucket)
	errorsCh := nonBlockingClient.Errors()
	go func() {
		for err := range errorsCh {
			err = errors.Wrap(err, "error occurred while writing metrics to influxdb")
			fn(err)
		}
	}()

	return &InfluxDb2Client{
		client: nonBlockingClient,
	}, nil
}
