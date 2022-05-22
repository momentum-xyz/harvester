package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

const (
	slotsPerSession = 600
	slotsPerEra     = 3600
)

type PalletStakingActiveEraInfo struct {
	Index types.U32       `json:"index"`
	Start types.OptionU64 `json:"start"`
}

type ChainSessionState struct {
	CurrentSessionIndex    types.U32 `json:"currentSessionIndex"`
	CurrentEra             uint32    `json:"currentEra"`
	SessionsPerEra         types.U32 `json:"sessionsPerEra"`
	TotalRewardPointsInEra types.U32 `json:"totalRewardPointsInEra"`
	SlotsPerSession        uint      `json:"slotsPerSession"`
	SlotsPerEra            uint      `json:"slotsPerEra"`
	ActiveEraStart         types.U64 `json:"activeEraStart"`
}

func (sh *SubstrateHarvester) ProcessChainSessionState(ctx context.Context, fn harvester.ErrorHandler) error {

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.publishChainSessionState()
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) publishChainSessionState() error {
	currentSessionIndexKey, err := sh.GetStorageDataKey("Session", "CurrentIndex")
	if err != nil {
		return errors.Wrap(err, " error while creating session current index storage data key")
	}
	var currentSessionIndex types.U32
	err = sh.GetStorageLatest(currentSessionIndexKey, &currentSessionIndex)
	if err != nil {
		return errors.Wrap(err, " error while fetching session current index")
	}

	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return errors.Wrap(err, " error while fetching active era for chain session state")
	}

	activeEraDepth, err := sh.GetEraDepth(activeEra)

	if err != nil {
		return errors.Wrap(err, "error while fetching active era depth")
	}

	totalRewardPointsInEra, err := sh.GetEraReward(activeEraDepth)

	if err != nil {
		return errors.Wrap(err, "error while fetching active era reward points")
	}

	if err != nil {
		return errors.Wrap(err, "error while fetching active era depth")
	}

	var sessionsPerEra types.U32

	err = sh.QueryConstant("Staking", "SessionsPerEra", &sessionsPerEra)

	if err != nil {
		return errors.Wrap(err, "error while fetching session per era constant value")
	}

	activeEraInfoKey, err := sh.GetStorageDataKey("Staking", "ActiveEra")
	if err != nil {
		return errors.Wrap(err, " error while creating active era key storage data key")
	}

	var activeEraInfo PalletStakingActiveEraInfo

	err = sh.GetStorageLatest(activeEraInfoKey, &activeEraInfo)
	if err != nil {
		return errors.Wrap(err, " error while fetching active era info")
	}

	exists, start := activeEraInfo.Start.Unwrap()
	if !exists {
		return errors.New("no value exists on active era unwarp")
	}

	sessionState := ChainSessionState{
		CurrentSessionIndex:    currentSessionIndex,
		CurrentEra:             activeEra,
		TotalRewardPointsInEra: totalRewardPointsInEra.Total,
		SlotsPerSession:        slotsPerSession,
		SlotsPerEra:            slotsPerEra,
		SessionsPerEra:         sessionsPerEra,
		ActiveEraStart:         start,
	}

	sessionStateJson, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}

	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/session", sh.cfg.Name), string(sessionStateJson))
	if err != nil {
		return err
	}
	return nil
}
