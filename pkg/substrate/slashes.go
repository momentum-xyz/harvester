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

type ValidatorSlashInEra struct {
	Perbill   types.U32
	BalanceOf types.U128
}

type Slash struct {
	AccountID    string `json:"accountId"`
	Amount       int64  `json:"amount"`
	EraIndex     uint32 `json:"eraIndex"`
	SessionIndex uint32 `json:"sessionIndex"`
	Type         string `json:"type"`
}

func (sh *SubstrateHarvester) ProcessSlashes(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	return sh.processSlashes(ctx, fn, pmc, topic, 5*time.Minute)
}

func (sh *SubstrateHarvester) processSlashes(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string,
	d time.Duration) error {

	log.Debug("processing slashes")

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.getSlashes(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) getSlashes(fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)

	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return errors.Wrap(err, "error while fetching active era")
	}

	sessionIndex, err := sh.GetSessionIndex()
	if err != nil {
		return errors.Wrap(err, "error while fetching current session index")
	}

	validatorAccountIDs, err := sh.getCurrentSessionValidators()
	if err != nil {
		return errors.Wrap(err, "error while fetching current session validators")
	}

	var validatorAddresses []string
	for _, accountID := range validatorAccountIDs {
		address, err := AccountIdToString(accountID, sh.cfg.Name)
		if err != nil {
			fn(err)
			continue
		}
		validatorAddresses = append(validatorAddresses, address)
	}

	slashes, err := sh.getValidatorSlashesInEra(fn, activeEra)
	if err != nil {
		return errors.Wrap(err, "error while fetching slashes")
	}

	for _, slash := range slashes {
		if Contains(validatorAddresses, slash.AccountID) {
			slash.SessionIndex = uint32(sessionIndex)
			err = sh.publishSlashEvent(slash, topic)
			if err != nil {
				fn(err)
			}
		}
	}

	return nil
}

func (sh *SubstrateHarvester) getValidatorSlashesInEra(fn harvester.ErrorHandler, era uint32) ([]Slash, error) {
	result := []Slash{}
	eraDepth, _ := sh.GetEraDepth(era)
	key, err := sh.CreateStorageKeyUnsafe("Staking", "ValidatorSlashInEra", eraDepth)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating validator slash in era storage data key")
	}

	keys, err := sh.GetKeysLatest(key)
	if err != nil {
		return nil, err
	}

	changes, err := sh.QueryStorageAt(keys)
	if err != nil {
		return nil, err
	}

	for _, change := range changes {
		var r ValidatorSlashInEra
		err = types.DecodeFromBytes(change.StorageData, &r)
		if err != nil {
			fn(err)
			continue
		}

		accountId := types.NewAccountID(change.StorageKey[len(change.StorageKey)-32:])
		address, err := AccountIdToString(accountId, sh.cfg.Name)
		if err != nil {
			fn(err)
			continue
		}

		result = append(result, Slash{AccountID: address, Amount: r.BalanceOf.Int64(), EraIndex: era, Type: "validator"})
	}

	return result, nil
}

func (sh *SubstrateHarvester) publishSlashEvent(slashEvent Slash, topic string) error {
	slashJson, err := json.Marshal(slashEvent)
	if err != nil {
		return err
	}

	account, era := slashEvent.AccountID, slashEvent.EraIndex
	log.Debug(fmt.Sprintf("%s - Publishing slash event for account: %s in era: %d", sh.cfg.Name, account, era))
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/%s/%s", sh.cfg.Name, topic, account), string(slashJson))
	if err != nil {
		return err
	}
	return nil
}
