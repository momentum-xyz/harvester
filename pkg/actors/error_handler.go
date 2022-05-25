package actors

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"

	"github.com/getsentry/sentry-go"
)

func ActorErrorHandler(err error) {
	log.Error(err)
	hub := sentry.CurrentHub().Clone()
	hub.Recover(err)
}
