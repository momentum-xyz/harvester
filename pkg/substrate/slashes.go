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
	Perbill types.U32
	Amount  types.U128
}

type Slash struct {
	AccountID    string `json:"accountId"`
	Amount       int64  `json:"amount"`
	EraIndex     uint32 `json:"eraIndex"`
	SessionIndex uint32 `json:"sessionIndex"`
}

func (sh *SubstrateHarvester) ProcessSlashes(ctx context.Context, fn harvester.ErrorHandler) error {
	log.Debug("processing slashes")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.getSlashes()
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) getSlashes() error {
	var slashes []Slash
	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return errors.Wrap(err, "error while fetching active era")
	}

	sessionIndex, err := sh.GetSessionIndex()
	if err != nil {
		return errors.Wrap(err, "error while fetching current session index")
	}

	activeValidators, err := sh.getCurrentSessionValidators()
	if err != nil {
		return errors.Wrap(err, "error while fetching current session validators")
	}

	for _, accountID := range activeValidators {
		address, err := accountIdToString(accountID)
		if err != nil {
			return err
		}

		eraDepth, _ := sh.GetEraDepth(activeEra)
		key, err := sh.GetStorageDataKey("Staking", "ValidatorSlashInEra", eraDepth, accountID[:])
		if err != nil {
			return errors.Wrap(err, "error while creating validator slash in era storage data key")
		}

		var amount types.U128
		err = sh.GetStorageLatest(key, &amount)
		if err != nil {
			return errors.Wrapf(err, "error while fetching validator slash in era:%v and accountId:", activeEra, address)
		}

		if amount.Int != nil {
			slashes = append(slashes, Slash{
				AccountID:    address,
				Amount:       amount.Int64(),
				EraIndex:     activeEra,
				SessionIndex: uint32(sessionIndex),
			})
		}
	}

	if len(slashes) > 0 {
		slashesJson, err := json.Marshal(slashes)
		if err != nil {
			return err
		}

		log.Logln(0, fmt.Sprintf("%s - Publishing slash event for era %d", sh.cfg.Name, activeEra))
		err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/slashes", sh.cfg.Name), string(slashesJson))
		if err != nil {
			return err
		}
	}
	return nil
}
