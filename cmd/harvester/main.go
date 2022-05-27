package main

import (
	"os"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/actors"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/getsentry/sentry-go"
)

func main() {

	var err error

	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = "config.yaml"
	}
	cfg := harvester.GetConfig(configPath, true)
	cfg.PrettyPrint()

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.Dsn,
		Environment:      cfg.Sentry.Environment,
		Release:          cfg.Sentry.Release,
		Debug:            cfg.Sentry.DebugEnable,
		AttachStacktrace: cfg.Sentry.AttachStacktrace,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return event

		},
	})

	if err != nil {
		log.Errorf("sentry.Init: %s", err)
		panic(err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(5 * time.Second)

	log.DefaultLogger, err = log.NewLogger(cfg.LogLevel.Level)

	defer log.DefaultLogger.Sync()

	if err != nil {
		sentry.CurrentHub().Recover(err)
		panic(err)
	}

	actors, err := actors.Start(cfg, actors.ActorErrorHandler)

	if err != nil {
		log.Error(err)
		sentry.CurrentHub().Recover(err)
		panic(err)
	}

	log.Info("Harvester actors started")
	log.Debug(actors.Run())
	log.Info("Harvester actors terminated")

}
