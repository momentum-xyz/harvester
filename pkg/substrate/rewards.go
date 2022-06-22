package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey"
)

type Reward struct {
	Era         uint32            `json:"era"`
	TotalPoints uint32            `json:"totalPoints"`
	Rewards     map[string]uint32 `json:"rewards"`
}

type EraRewardPoints struct {
	Total      types.U32
	Individual []struct {
		AccountID    types.AccountID
		RewardPoints types.U32
	}
}

func (sh *SubstrateHarvester) ProcessErasRewardPoints(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			sh.publishErasRewardPoints(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}

}

func (sh *SubstrateHarvester) publishErasRewardPoints(fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)

	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return errors.Wrap(err, "error while fetching active era")
	}

	activeEraDepth, err := sh.GetEraDepth(activeEra)

	if err != nil {
		return errors.Wrap(err, "error while fetching active era depth")
	}

	rewards, err := sh.GetEraReward(activeEraDepth)
	if err != nil {
		return err
	}

	rewardsMap := make(map[string]uint32)
	for _, v := range rewards.Individual {
		address, err := subkey.SS58Address(v.AccountID[:], 2)
		if err != nil {
			return err
		}

		rewardsMap[address] = uint32(v.RewardPoints)
	}

	rewardsJson, err := json.Marshal(Reward{
		Era:         activeEra,
		TotalPoints: uint32(rewards.Total),
		Rewards:     rewardsMap,
	})
	if err != nil {
		return err
	}

	log.Logln(0, fmt.Sprintf("%s - Publishing reward event for era %d", sh.cfg.Name, activeEra))
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/%s", sh.cfg.Name, topic), string(rewardsJson))
	if err != nil {
		return err
	}
	return nil
}

func (sh *SubstrateHarvester) GetEraReward(era []byte) (EraRewardPoints, error) {

	noRewards := EraRewardPoints{}

	key, err := sh.GetStorageDataKey("Staking", "ErasRewardPoints", era)
	if err != nil {
		return noRewards, errors.Wrapf(err, "error while fetching storage key for erasRewardPoints for era %v", spew.Sdump(era))
	}

	var rewards EraRewardPoints
	err = sh.GetStorageLatest(key, &rewards)
	if err != nil {
		return noRewards, errors.Wrapf(err, "error while fetching storage erasRewardPoints for key %v", spew.Sdump(key))
	}

	return rewards, nil
}
