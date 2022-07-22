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
	ActiveNominators      int        `json:"activeNominators"`
	CandidateValidators   uint32     `json:"candidateValidators"`
	TotalStakeInActiveEra types.U128 `json:"totalStakeInActiveEra"`
	TotalStakeInLastEra   types.U128 `json:"totalStakeInLastEra"`
	LastEraReward         types.U128 `json:"lastEraReward"`
	StakingRatio          float32    `json:"stakingRatio"`
	IdealStakingRatio     float32    `json:"idealStakingRatio"`
	StakedReturn          float32    `json:"stakedReturn"`
	Inflation             float32    `json:"inflation"`
}

func (sh *SubstrateHarvester) ProcessChainEraState(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	return sh.processChainEraState(ctx, fn, pmc, topic, 5*time.Minute)
}

func (sh *SubstrateHarvester) processChainEraState(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string,
	d time.Duration) error {

	ticker := time.NewTicker(d)

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
		return errors.Wrap(err, " error while fetching activeValidators")
	}

	activeNominators, err := sh.getActiveNominators(fn)
	if err != nil {
		return errors.Wrap(err, " error while fetching activeNominators")
	}

	stakingRatio, err := sh.getStakingRatio(activeEra)
	if err != nil {
		return err
	}

	inflation, err := sh.GetInflation(stakingRatio)
	if err != nil {
		return errors.Wrap(err, " error while fetching inflation")
	}

	eraState := ChainEraState{
		Name:                  sh.cfg.Name,
		ActiveEra:             activeEra,
		ActiveValidators:      len(activeValidators),
		ActiveNominators:      len(activeNominators),
		TotalStakeInActiveEra: totalStakeInActiveEra,
		TotalStakeInLastEra:   totalStakeInLastEra,
		LastEraReward:         reward,
		StakingRatio:          stakingRatio,
		IdealStakingRatio:     inflation.IdealStake,
		StakedReturn:          inflation.StakedReturn,
		Inflation:             inflation.Inflation,
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

	stakingRatio := (float64(totalStaked.Uint64()) * 100) / float64(totalIssuance.Uint64())
	return float32(stakingRatio), nil
}
