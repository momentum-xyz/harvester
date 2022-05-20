package wire

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mysql"
	"github.com/OdysseyMomentumExperience/harvester/pkg/publisher"
	"github.com/OdysseyMomentumExperience/harvester/pkg/repository"
)

func NewHarvester(cfg *harvester.Config, fn harvester.ErrorHandler) (*harvester.Harvester, func(), error) {
	db, cleanupDB, err := mysql.NewDB(&cfg.MySQL)
	if err != nil {
		return nil, nil, err
	}

	repository, err := repository.NewRepository(db, cfg.MySQL.Migrate)
	if err != nil {
		return nil, nil, err
	}

	mqttClient := mqtt.GetMQTTClient(&cfg.MQTT, fn)

	publisher, err := publisher.NewPublisher(mqttClient)
	if err != nil {
		return nil, nil, err
	}

	harvester, err := harvester.NewHarvester(cfg, repository, publisher)
	if err != nil {
		cleanupDB()
		return nil, nil, err
	}

	cleanup := func() { cleanupDB() }

	return harvester, cleanup, nil
}
