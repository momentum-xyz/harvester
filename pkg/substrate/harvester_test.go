package substrate

import (
	"reflect"
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mysql"
	performancemonitor "github.com/OdysseyMomentumExperience/harvester/pkg/performance_monitor"
	"github.com/OdysseyMomentumExperience/harvester/pkg/publisher"
	"github.com/OdysseyMomentumExperience/harvester/pkg/repository"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var mockCfg = harvester.Config{
	MQTT: mqtt.Config{
		Host:     "localhost",
		Port:     1883,
		ClientId: "harvester",
	},
	MySQL: mysql.Config{
		Database: "harvester",
		Host:     "localhost",
		Password: "",
		Username: "root",
		Port:     3306,
	},
	InfluxDB: harvester.InfluxDbConfig{
		Url:                   "http://localhost:8086",
		Org:                   "",
		Bucket:                "",
		Token:                 "",
		BatchSize:             5,
		RetryInterval:         2,
		SystemMetricsInterval: 10,
	},
}
var mockChainCfg = harvester.ChainConfig{
	Name:         "harvester",
	RPC:          "ws://localhost:9944",
	ActiveTopics: []string{"block-creation-event"},
}

var mockPmc, _ = performancemonitor.NewPerformanceMonitorClient(&mockCfg, func(err error) {})
var mockDB, _, _ = mysql.NewDB(&mockCfg.MySQL)
var mockRepository, _ = repository.NewRepository(mockDB, mockCfg.MySQL.Migrate)
var mqttClient = mqtt.GetMQTTClient(&mockCfg.MQTT, func(err error) {})
var mockPublisher, _ = publisher.NewPublisher(mqttClient)
var mockHarvester, _ = harvester.NewHarvester(&mockCfg, mockRepository, mockPublisher, mockPmc)
var mockSh, _ = NewHarvester(mockChainCfg, mockHarvester.Publisher, mockHarvester.Repository)

func TestSubstrateHarvester(t *testing.T) {
	t.Run("NewHarvester()", func(t *testing.T) {
		_mockSh, err := NewHarvester(mockChainCfg, mockHarvester.Publisher, mockHarvester.Repository)
		assert.Nil(t, err)
		assert.IsType(t, _mockSh, &SubstrateHarvester{})

		_mockChainCfg := mockChainCfg
		_mockChainCfg.RPC = "ws://example.com"
		_mockSh, err = NewHarvester(_mockChainCfg, mockHarvester.Publisher, mockHarvester.Repository)
		assert.NotNil(t, err)
		assert.Nil(t, _mockSh)
	})

	t.Run("getApi()", func(t *testing.T) {
		api, err := getApi(mockChainCfg.RPC)
		assert.Nil(t, err)
		assert.IsType(t, api, &gsrpc.SubstrateAPI{})

		api, err = getApi("ws://example.com")
		assert.Nil(t, api)
		assert.NotNil(t, err)
	})

	t.Run("getMetadata()", func(t *testing.T) {
		api, err := getApi(mockChainCfg.RPC)
		assert.Nil(t, err)
		assert.IsType(t, api, &gsrpc.SubstrateAPI{})

		metadata, err := getMetadata(api)
		assert.IsType(t, metadata, &types.Metadata{})
		assert.Nil(t, err)
	})

	t.Run("getNetworkID()", func(t *testing.T) {
		assert.Equal(t, getNetworkID("kusama"), uint8(2))
		assert.Equal(t, getNetworkID("polkadot"), uint8(1))
		assert.Equal(t, getNetworkID("aaaa"), uint8(2))
	})

	t.Run("getActiveProcesses()", func(t *testing.T) {
		processes := mockSh.getActiveProcesses()
		assert.Equal(t, len(processes), 1)
		assert.Equal(t, processes[0].Topic, "block-creation-event")
		assert.Equal(t, reflect.ValueOf(processes[0].Process).Pointer(), reflect.ValueOf(mockSh.ProcessNewHeader).Pointer())
	})

	t.Run("topicProcessorStore()", func(t *testing.T) {
		process := mockSh.topicProcessorStore()("block-creation-event")
		assert.Equal(t, reflect.ValueOf(process).Pointer(), reflect.ValueOf(mockSh.ProcessNewHeader).Pointer())

		process = mockSh.topicProcessorStore()("aaaaa")
		assert.Nil(t, process)
	})

	t.Run("Stop()", func(t *testing.T) {
		mockSh.Stop()
	})
}
