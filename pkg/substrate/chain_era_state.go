package substrate

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type ChainEraState struct {
	Name                  string     `json:"name"`
	ActiveEra             uint32     `json:"activeEra"`
	ActiveValidators      uint32     `json:"activeValidators"`
	CandidateValidators   uint32     `json:"candidateValidators"`
	TotalStakeInActiveEra types.U128 `json:"totalStakeInActiveEra"`
	TotalStakeInLastEra   types.U128 `json:"totalStakeInLastEra"`
	LastEraReward         types.U128 `json:"lastEraReward"`
}

func (sh *SubstrateHarvester) GetChainEraState(activeEra uint32,
	activeValidators uint32,
	candidateValidators uint32,
	fn harvester.ErrorHandler) error {

	totalStakeInActiveEra, err := sh.getEraTotalStake(activeEra)
	if err != nil {
		fn(err)
		return err
	}

	totalStakeInLastEra, err := sh.getEraTotalStake(activeEra - 1)
	if err != nil {
		fn(err)
		return err
	}

	reward, err := sh.GetErasValidatorReward(activeEra - 1)
	if err != nil {
		fn(err)
		return err
	}
	eraState := ChainEraState{
		Name:                  sh.cfg.Name,
		ActiveEra:             activeEra,
		ActiveValidators:      activeValidators,
		CandidateValidators:   candidateValidators,
		TotalStakeInActiveEra: totalStakeInActiveEra,
		TotalStakeInLastEra:   totalStakeInLastEra,
		LastEraReward:         reward,
	}

	msg, err := json.Marshal(eraState)

	if err != nil {
		fn(err)
		return err
	}
	log.Infof("%s - Publishing Chain Era State For Era %d", sh.cfg.Name, activeEra)
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/era", sh.cfg.Name), string(msg))
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
