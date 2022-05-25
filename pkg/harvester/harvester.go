package harvester

import "context"

type ErrorHandler func(error)

type Harvester struct {
	Cfg        *Config
	Repository Repository
	Publisher  Publisher
	// map of chain harvesters
}

type ActiveHarvesterProcess func(context.Context, ErrorHandler) error

func NewHarvester(cfg *Config, repository Repository, publisher Publisher) (*Harvester, error) {
	return &Harvester{
		Cfg:        cfg,
		Repository: repository,
		Publisher:  publisher,
	}, nil
}
