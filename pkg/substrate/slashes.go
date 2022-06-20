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

func (sh *SubstrateHarvester) ProcessSlashes(ctx context.Context, fn harvester.ErrorHandler) error {
	return sh.processSlashes(ctx, fn, 5*time.Minute)
}

func (sh *SubstrateHarvester) processSlashes(ctx context.Context, fn harvester.ErrorHandler, d time.Duration) error {
	log.Debug("processing slashes")

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.getSlashes(fn)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) getSlashes(fn harvester.ErrorHandler) error {
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

	for _, accountID := range validatorAccountIDs {
		address, err := accountIdToString(accountID)
		if err != nil {
			fn(err)
			continue
		}

		validatorSlashInEra, err := sh.getValidatorSlashInEra(activeEra, address)
		if err != nil {
			fn(err)
			continue
		}

		if validatorSlashInEra.BalanceOf.Int != nil {
			err = sh.publishSlashEvent(Slash{
				AccountID:    address,
				Amount:       validatorSlashInEra.BalanceOf.Int64(),
				EraIndex:     activeEra,
				SessionIndex: uint32(sessionIndex),
				Type:         "validator",
			})
			if err != nil {
				fn(err)
			}
		}
	}

	return nil
}

func (sh *SubstrateHarvester) getValidatorSlashInEra(era uint32, address string) (*ValidatorSlashInEra, error) {
	var validatorSlashInEra ValidatorSlashInEra
	eraDepth, _ := sh.GetEraDepth(era)
	accountID, err := stringToAccountId(address)

	if err != nil {
		return nil, err
	}

	key, err := sh.GetStorageDataKey("Staking", "ValidatorSlashInEra", eraDepth, accountID[:])
	if err != nil {
		return nil, errors.Wrap(err, "error while creating validator slash in era storage data key")
	}

	err = sh.GetStorageLatest(key, &validatorSlashInEra)
	if err != nil {
		return nil, errors.Wrapf(err, "error while fetching validator slash in era:%d and accountId:%s", era, address)
	}

	return &validatorSlashInEra, nil
}

func (sh *SubstrateHarvester) publishSlashEvent(slashEvent Slash) error {
	slashJson, err := json.Marshal(slashEvent)
	if err != nil {
		return err
	}

	account, era := slashEvent.AccountID, slashEvent.EraIndex
	log.Logln(0, fmt.Sprintf("%s - Publishing slash event for account:%s in era %d", sh.cfg.Name, account, era))
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/slashes/%s", sh.cfg.Name, account), string(slashJson))
	if err != nil {
		return err
	}
	return nil
}
