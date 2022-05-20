package mqtt

// Config : structure to hold MQTT configuration
type Config struct {
	Host                 string `yaml:"host" envconfig:"MQTT_BROKER_HOST"`
	Port                 uint   `yaml:"port" envconfig:"MQTT_BROKER_PORT"`
	User                 string `yaml:"user" envconfig:"MQTT_BROKER_USER"`
	Password             string `yaml:"password" envconfig:"MQTT_BROKER_PASSWORD"`
	ClientId             string `yaml:"clientId" envconfig:"MQTT_CLIENTID"`
	Topics               topics `yaml:"topics" envconfig:"MQTT_TOPICS"`
	MaxReconnectAttempts uint   `yaml:"maxReconnectAttempts" envconfig:"MQTT_MAX_RECONNECT_ATTEMPTS"`
	ReconnectInterval    uint   `yaml:"reconnectInterval" envconfig:"MQTT_RECONNECT_INTERVAL"`
}

type topics struct {
	ExampleTopic string `yaml:"exampleTopic" envconfig:"MQTT_TOPICS_EXAMPLE_TOPIC"`
}

func (x *Config) Init() {
	x.Host = "localhost"
	x.Port = 1883
	x.User = ""
	x.Password = ""
	x.MaxReconnectAttempts = 3
	x.ReconnectInterval = 3
}
