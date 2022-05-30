package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
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
	SlotsPerSession        types.U64 `json:"slotsPerSession"`
	SlotsPerEra            types.U64 `json:"slotsPerEra"`
	SlotDuration           types.U64 `json:"slotDuration"`
	ActiveEraStart         types.U64 `json:"activeEraStart"`
	CurrentSlotInSession   types.U64 `json:"currentSlotInSession"`
	CurrentSlotInEra       types.U32 `json:"currentSlotInEra"`
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

	slotDuration, err := sh.GetSlotDuration()

	if err != nil {
		return err
	}

	sessionLength, err := sh.GetSessionLength()

	if err != nil {
		return err
	}

	sessionsPerEra, err := sh.GetSessionPerEra()
	if err != nil {
		return err
	}

	currentSlotInSession, err := sh.currentSlotInSession(sessionLength)

	if err != nil {
		return err
	}

	currentSlotInEra, err := sh.currentSlotInEra(currentSlotInSession)
	if err != nil {
		return err
	}

	slotsPerEra := sessionLength * types.U64(sessionsPerEra)

	sessionState := ChainSessionState{
		CurrentSessionIndex:    currentSessionIndex,
		CurrentEra:             activeEra,
		TotalRewardPointsInEra: totalRewardPointsInEra.Total,
		CurrentSlotInSession:   currentSlotInSession,
		SlotsPerSession:        sessionLength,
		CurrentSlotInEra:       currentSlotInEra,
		SlotsPerEra:            slotsPerEra,
		SessionsPerEra:         sessionsPerEra,
		ActiveEraStart:         start,
		SlotDuration:           slotDuration,
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

func (sh *SubstrateHarvester) currentSlotInSession(sessionLength types.U64) (types.U64, error) {

	epochIndexKey, err := sh.GetStorageDataKey("Babe", "EpochIndex")
	if err != nil {
		return 0, errors.Wrap(err, " error while creating babe epoch index storage data key")
	}
	var epochIndex types.U64

	err = sh.GetStorageLatest(epochIndexKey, &epochIndex)
	if err != nil {
		return 0, errors.Wrap(err, " error while fetching epoch index")
	}

	epochOrGenesisStartSlotKey, err := sh.GetStorageDataKey("Babe", "GenesisSlot")
	if err != nil {
		return 0, errors.Wrap(err, " error while creating babe genesisSlot storage data key")
	}
	var epochOrGenesisStartSlot types.U64

	err = sh.GetStorageLatest(epochOrGenesisStartSlotKey, &epochOrGenesisStartSlot)
	if err != nil {
		return 0, errors.Wrap(err, " error while fetching epoch genesis slot")
	}

	epochStartSlot := (epochIndex * sessionLength) + epochOrGenesisStartSlot

	currentSlotKey, err := sh.GetStorageDataKey("Babe", "CurrentSlot")
	if err != nil {
		return 0, errors.Wrap(err, " error while creating babe current slot storage data key")
	}
	var currentSlot types.U64

	err = sh.GetStorageLatest(currentSlotKey, &currentSlot)
	if err != nil {
		return 0, errors.Wrap(err, " error while fetching current slot")
	}

	sessionProgress := currentSlot - epochStartSlot

	return sessionProgress, nil

}

func (sh *SubstrateHarvester) currentSlotInEra(sessionProgress types.U64) (types.U32, error) {

	sessionCurrentIndexKey, err := sh.GetStorageDataKey("Session", "CurrentIndex")
	if err != nil {
		return 0, errors.Wrap(err, " error while creating session current index storage data key")
	}
	var currentIndex types.U32

	err = sh.GetStorageLatest(sessionCurrentIndexKey, &currentIndex)
	if err != nil {
		return 0, errors.Wrap(err, " error while fetching session current index")
	}

	activeEra, err := sh.GetActiveEraDepth()

	if err != nil {
		return 0, errors.Wrap(err, " error while fetching active era for erasStartSessionIndex")
	}

	activeEraStartSessionIndexKey, err := sh.GetStorageDataKey("Staking", "ErasStartSessionIndex", activeEra)
	if err != nil {
		return 0, errors.Wrap(err, " error while creating session current index storage data key")
	}
	var activeEraStartSessionIndex types.U32

	err = sh.GetStorageLatest(activeEraStartSessionIndexKey, &activeEraStartSessionIndex)
	if err != nil {
		return 0, errors.Wrap(err, " error while fetching session current index")
	}

	sessionLength, err := sh.GetSessionLength()
	if err != nil {
		return 0, err
	}
	log.Debugf(" current slot in era calucaltion sessionLength  %v, sessionProgress %v ", sessionLength, sessionProgress)

	eraProgress := (currentIndex-activeEraStartSessionIndex)*types.U32(sessionLength) + types.U32(sessionProgress)

	log.Debugf(" current slot in era %v", eraProgress)

	return eraProgress, nil
}
