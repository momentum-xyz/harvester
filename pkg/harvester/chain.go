package harvester

import (
	"context"
)

// Chain harvesters should implement this interface
type ChainHarvester interface {
	Start(ctx context.Context, pmc PerformanceMonitorClient, fn ErrorHandler) error
	Stop()
}

type ChainConfig struct {
	Name         string   `yaml:"name"`
	RPC          string   `yaml:"rpc"`
	Type         string   `yaml:"type"`
	FromBlock    uint64   `yaml:"fromBlock"`
	ActiveTopics []string `yaml:"active_topics" envconfig:"ACTIVE_TOPICS"`
}
