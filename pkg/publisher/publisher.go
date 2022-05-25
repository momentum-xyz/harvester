package publisher

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
)

type Publisher struct {
	Mqtt *mqtt.Mqtt
}

func NewPublisher(mqtt *mqtt.Mqtt) (harvester.Publisher, error) {
	return &Publisher{
		Mqtt: mqtt,
	}, nil
}

func (p *Publisher) Publish(topic string, message string) error {
	return p.Mqtt.Publish(topic, 1, message, false)
}

func (p *Publisher) PublishRetained(topic string, message string) error {
	return p.Mqtt.Publish(topic, 1, message, true)
}
