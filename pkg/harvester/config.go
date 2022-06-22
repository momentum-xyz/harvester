package harvester

import (
	"fmt"

	"github.com/pborman/getopt/v2"
	"go.uber.org/zap/zapcore"

	"io"
	"os"

	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mysql"
	"github.com/OdysseyMomentumExperience/harvester/pkg/sentry"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type LogLevel struct {
	zapcore.Level
}

func (l *LogLevel) UnmarshalFlag(level string) error {
	return l.Set(level)
}

// Config : structure to hold configuration
type Config struct {
	MQTT                 mqtt.Config                `yaml:"mqtt"`
	MySQL                mysql.Config               `yaml:"mysql"`
	LogLevel             *LogLevel                  `yaml:"loglevel" env:"LOG_LEVEL"`
	Chains               []ChainConfig              `yaml:"chains"`
	EnabledChains        []string                   `yaml:"enabled_chains" envconfig:"ENABLED_CHAINS"`
	ExchangeRateProvider ExchangeRateProviderConfig `yaml:"exchange_rate_provider"`
	Sentry               sentry.Config              `yaml:"sentry" envconfig:"SENTRY"`
	InfluxDB             InfluxDbConfig             `yaml:"influx_db" envconfig:"INFLUX_DB"`
}

func (x *Config) Init() {
	x.MQTT.Init()
	x.MySQL.Init()
}

func (cfg *Config) defConfig() {
	cfg.Init()
}

func (cfg *Config) readOpts() {
	helpFlag := false
	getopt.Flag(&helpFlag, 'h', "display help")

	getopt.Parse()
	if helpFlag {
		getopt.Usage()
		os.Exit(0)
	}
}

func (cfg *Config) processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func (cfg *Config) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (cfg *Config) readFile(path string) {
	if !cfg.fileExists(path) {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		cfg.processError(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		if err != io.EOF {
			cfg.processError(err)
		}
	}
}

func (cfg *Config) readEnv() {
	err := envconfig.Process("", cfg)
	if err != nil {
		cfg.processError(err)
	}
}

func (cfg *Config) PrettyPrint() {
	d, _ := yaml.Marshal(cfg)
	log.Logf(1, "--- Config ---\n%s\n\n", string(d))
}

// GetConfig : get config file
func GetConfig(path string, enableFlags bool) Config {
	var cfg Config
	cfg.defConfig()

	cfg.readFile(path)
	cfg.readEnv()

	if enableFlags {
		cfg.readOpts()
	}
	cfg.PrettyPrint()
	return cfg
}
