package substrate

import (
	"context"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/oklog/run"
)

func (sh *SubstrateHarvester) Start(ctx context.Context, fn harvester.ErrorHandler) error {
	log.Infof("starting chain harvester %s\n", sh.cfg.Name)

	// Start substrate process
	g := new(run.Group)

	for _, substrateProcess := range sh.getActiveProcesses() {

		{
			substrateProcess := substrateProcess
			ctx, cancel := context.WithCancel(ctx)
			g.Add(func() error {
				return substrateProcess(ctx, fn)
			}, func(error) {
				cancel()
			})
		}
	}

	{

		g.Add(func() error {
			<-ctx.Done()
			return ctx.Err()
		}, func(err error) {
		})

	}

	return g.Run()
}
