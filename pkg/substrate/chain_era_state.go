package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

type ChainEraState struct {
	Name                  string     `json:"name"`
	ActiveEra             uint32     `json:"activeEra"`
	ActiveValidators      int        `json:"activeValidators"`
	CandidateValidators   uint32     `json:"candidateValidators"`
	TotalStakeInActiveEra types.U128 `json:"totalStakeInActiveEra"`
	TotalStakeInLastEra   types.U128 `json:"totalStakeInLastEra"`
	LastEraReward         types.U128 `json:"lastEraReward"`
	StakingRatio          float32    `json:"stakingRatio"`
}

func (sh *SubstrateHarvester) ProcessChainEraState(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	ticker := time.NewTicker(5 * time.Minute)

	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.publishEraInfo(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}

}

func (sh *SubstrateHarvester) publishEraInfo(fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)

	activeEra, err := sh.GetActiveEra()

	if err != nil {
		return errors.Wrap(err, " error while fetching active era")
	}

	totalStakeInActiveEra, err := sh.getEraTotalStake(activeEra)
	if err != nil {
		return errors.Wrap(err, " error while fetching totalStakeInActiveEra")
	}

	totalStakeInLastEra, err := sh.getEraTotalStake(activeEra - 1)
	if err != nil {
		return errors.Wrap(err, " error while fetching totalStakeInLastEra")
	}

	reward, err := sh.GetErasValidatorReward(activeEra - 1)
	if err != nil {

		return err
	}

	activeValidators, err := sh.getCurrentSessionValidators()
	if err != nil {
		return err
	}

	stakingRatio, err := sh.getStakingRatio(activeEra)
	if err != nil {
		return err
	}

	eraState := ChainEraState{
		Name:                  sh.cfg.Name,
		ActiveEra:             activeEra,
		ActiveValidators:      len(activeValidators),
		TotalStakeInActiveEra: totalStakeInActiveEra,
		TotalStakeInLastEra:   totalStakeInLastEra,
		LastEraReward:         reward,
		StakingRatio:          stakingRatio,
	}

	msg, err := json.Marshal(eraState)

	if err != nil {
		fn(err)
		return err
	}
	log.Infof("%s - Publishing Chain Era State For Era %d", sh.cfg.Name, activeEra)
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/%s", sh.cfg.Name, topic), string(msg))
	if err != nil {
		fn(err)
		return err
	}
	return nil
}

func (sh *SubstrateHarvester) GetErasValidatorReward(era uint32) (types.U128, error) {

	zeroReward := types.NewU128(*big.NewInt(0))

	historyDepth, err := sh.GetEraDepth(era)
	if err != nil {
		return zeroReward, err
	}

	erasValidatorRewardKey, err := sh.GetStorageDataKey("Staking", "ErasValidatorReward", historyDepth)

	if err != nil {
		return zeroReward, err
	}

	var erasValidatorReward types.U128
	err = sh.GetStorageLatest(erasValidatorRewardKey, &erasValidatorReward)

	if err != nil || erasValidatorReward.Int == nil {
		return zeroReward, err
	}

	return erasValidatorReward, nil
}

//TODO <nil> check
func (sh *SubstrateHarvester) getEraTotalStake(era uint32) (types.U128, error) {
	totalStakeInEra := types.NewU128(*big.NewInt(0))

	eraDepth, err := sh.GetEraDepth(era)

	if err != nil {
		return totalStakeInEra, err
	}

	totalStakeInEraKey, err := sh.GetStorageDataKey("Staking", "ErasTotalStake", eraDepth)
	if err != nil {
		return totalStakeInEra, err
	}

	err = sh.GetStorageLatest(totalStakeInEraKey, &totalStakeInEra)

	if err != nil {
		return totalStakeInEra, err
	}

	return totalStakeInEra, nil
}

func (sh *SubstrateHarvester) getStakingRatio(activeEra uint32) (float32, error) {
	totalIssuance, err := sh.GetTotalIssuance()
	if err != nil || totalIssuance.Int == nil {
		return 0, err
	}

	totalStaked, err := sh.getEraTotalStake(activeEra)
	if err != nil {
		return 0, err
	}

	stakingRatio := (float64(totalStaked.Uint64()) * float64(100)) / float64(totalIssuance.Uint64())
	return float32(stakingRatio), nil
}
