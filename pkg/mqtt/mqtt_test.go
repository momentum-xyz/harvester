package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg = Config{
	Host:     "localhost",
	Port:     1883,
	ClientId: "harvester",
}

func TestInit(t *testing.T) {
	_cfg := Config{}
	_cfg.Init()
	assert.Equal(t, _cfg.Host, "localhost")
	assert.Equal(t, _cfg.Port, uint(1883))
}

func TestSubscribe(t *testing.T) {
	t.Run("success subscribe", func(t *testing.T) {
		ch := make(chan []byte)
		mqttClient := GetMQTTClient(&cfg, func(err error) {})
		err := mqttClient.Subscribe("test_topic", 1, ch)
		assert.Nil(t, err)
	})

	t.Run("success subscribe", func(t *testing.T) {
		ch := make(chan []byte)
		mqttClient := GetMQTTClient(&Config{}, func(err error) {})
		err := mqttClient.Subscribe("test_topic", 1, ch)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Not Connected")
	})
}

func TestPublish(t *testing.T) {
	t.Run("success publish", func(t *testing.T) {
		mqttClient := GetMQTTClient(&cfg, func(err error) {})
		err := mqttClient.Publish("test_topic", 1, "test message", false)
		assert.Nil(t, err)
	})
	t.Run("publish on connection lost", func(t *testing.T) {
		_cfg := &Config{}
		mqttClient := GetMQTTClient(_cfg, func(err error) {})
		err := mqttClient.Publish("test_topic", 1, "test message", false)
		assert.Contains(t, err.Error(), "failed to publish message")
	})
	t.Run("handle messages on connection lost", func(t *testing.T) {
		mqttClient := GetMQTTClient(&cfg, func(err error) {})
		mqttClient.client.Disconnect(1)
		err := mqttClient.Publish("test_topic", 1, "test message", false)
		assert.Contains(t, err.Error(), "failed to publish message")
	})
}
