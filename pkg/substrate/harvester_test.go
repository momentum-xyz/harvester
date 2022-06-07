package substrate

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/wire"
)

var cfg = harvester.Config{MQTT: mqtt.Config{Host: "localhost", Port: 1883, ClientId: "harvester"}}
var chainCfg = harvester.ChainConfig{Name: "harvester", RPC: "ws://localhost:9944", ActiveTopics: []string{}}
var h, _, _ = wire.NewHarvester(&cfg, func(err error) {})
var sh, _ = NewHarvester(chainCfg, h.Publisher, h.Repository)
