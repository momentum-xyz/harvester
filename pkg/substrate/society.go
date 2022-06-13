package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
)

func (sh *SubstrateHarvester) ProcessSocietyMembers(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.getSocietyMembers(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()

		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) getSocietyMembers(fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)

	key, err := sh.GetStorageDataKey("Society", "Members")
	if err != nil {
		return err
	}

	var target []types.AccountID
	err = sh.GetStorageLatest(key, &target)
	if err != nil {
		return err
	}

	networkId := getNetworkID(sh.cfg.Name)
	var accountIDs []string
	for _, account := range target {
		accountID, err := subkey.SS58Address(account[:], networkId)
		if err != nil {
			return err
		}
		accountIDs = append(accountIDs, accountID)
	}

	accountsJson, err := json.Marshal(accountIDs)
	if err != nil {
		return err
	}

	err = sh.publisher.PublishRetained(fmt.Sprintf("harvester/%s/%s", sh.cfg.Name, topic), string(accountsJson))
	if err != nil {
		return err
	}
	return nil
}
