package harvester

type Publisher interface {
	Publish(topic string, message string) error
	PublishRetained(topic string, message string) error
}
