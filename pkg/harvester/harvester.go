package harvester

import (
	"context"
)

type ErrorHandler func(error)

type Harvester struct {
	Cfg                      *Config
	Repository               Repository
	Publisher                Publisher
	PerformanceMonitorClient PerformanceMonitorClient
	// map of chain harvesters
}

// pass first error handler
type ActiveHarvesterProcess func(context.Context, ErrorHandler, PerformanceMonitorClient, string) error

type TopicProcess struct {
	Topic   string
	Process ActiveHarvesterProcess
}

func NewHarvester(cfg *Config, repository Repository, publisher Publisher, pmc PerformanceMonitorClient) (*Harvester, error) {
	return &Harvester{
		Cfg:                      cfg,
		Repository:               repository,
		Publisher:                publisher,
		PerformanceMonitorClient: pmc,
	}, nil
}
