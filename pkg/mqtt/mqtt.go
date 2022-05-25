package mqtt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/pkg/errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTMessage struct {
	Topic    string
	Payload  string
	QOS      int
	Retained bool
}

type Mqtt struct {
	client mqtt.Client
}

var pendingMessages []MQTTMessage

func GetMQTTClient(cfg *Config, fn func(error)) *Mqtt {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + cfg.Host + ":" + strconv.FormatUint(uint64(cfg.Port), 10))
	opts.SetClientID(cfg.ClientId)
	opts.SetUsername(cfg.User)
	opts.SetPassword(cfg.Password)
	opts.SetMaxReconnectInterval(time.Duration(cfg.ReconnectInterval) * time.Second)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler(cfg, fn)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	return &Mqtt{client: client}
}

func (m *Mqtt) Subscribe(topic string, qos int, channel chan []byte) error {
	if token := m.client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		channel <- msg.Payload()
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *Mqtt) Publish(topic string, qos int, msg string, retained bool) error {
	if !m.client.IsConnectionOpen() {
		pendingMessages = append(pendingMessages, MQTTMessage{Topic: topic, Payload: msg})
		return errors.Wrap(errors.New("Connection to MQTT broker lost"), "failed to publish message")
	}
	if token := m.client.Publish(topic, byte(qos), retained, msg); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "failed to publish message")
	}
	return nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Infof("Connected to MQTT broker")

	go func() {
		n := 0
		for _, msg := range pendingMessages {
			token := client.Publish(msg.Topic, byte(msg.QOS), msg.Retained, msg.Payload)
			if token.Wait() && token.Error() != nil {
				pendingMessages[n] = msg
				n++
			}
		}
		pendingMessages = pendingMessages[:n]
	}()
}

func connectLostHandler(cfg *Config, fn func(error)) mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		log.Errorf("Connection to MQTT broker lost: %v\n", err)

		go func() {
			time.Sleep(time.Duration(cfg.ReconnectInterval) * time.Second)

			for i := 1; i <= int(cfg.MaxReconnectAttempts); i++ {
				time.Sleep(time.Duration(cfg.ReconnectInterval) * time.Second)

				if !client.IsConnectionOpen() && i == int(cfg.MaxReconnectAttempts) {
					err := errors.New(fmt.Sprintf("Connection to MQTT broker lost after %d attempts: %v\n", i, pendingMessages))
					fn(err)
					panic(err)
				} else if client.IsConnectionOpen() {
					break
				}
			}
		}()
	}
}
