package substrate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/log"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func (sh *SubstrateHarvester) ProcessValidators(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	return sh.processErasRewardPoints(ctx, fn, pmc, topic, 5*time.Minute)
}

func (sh *SubstrateHarvester) processValidators(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string,
	d time.Duration) error {
	log.Debug("processing validators")

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.updateValidators(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()
		}

		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) updateValidators(fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	// Get all known stashAccounts
	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)
	stashAccounts, err := sh.getStashAccounts()
	if err != nil {
		return err
	}
	// Get all active validators
	activeValidators, err := sh.getCurrentSessionValidators()
	if err != nil {
		return err
	}

	dbValidatorMap, err := sh.repo.GetValidatorMap()
	if err != nil {
		return err
	}

	for _, node := range stashAccounts {
		validatorNode, err := StringToAccountId(node)
		if err != nil {
			return err
		}
		validatorPrefs, err := sh.getValidatorPreferences(validatorNode)
		if err != nil {
			return err
		}
		commission := float64(validatorPrefs.Commission.Int64()) / 1000000000

		parentIdentity, err := sh.getSuperOf(validatorNode)
		if err != nil {
			return err
		}

		parentAccount, err := AccountIdToString(parentIdentity.Account, sh.cfg.Name)
		if err != nil {
			return err
		}

		subAccounts, err := sh.getSubsOf(validatorNode)
		if err != nil {
			return err
		}

		children := make([]string, 0)
		for _, c := range subAccounts.Accounts {
			childAccount, err := AccountIdToString(c, sh.cfg.Name)
			if err != nil {
				return err
			}
			children = append(children, childAccount)
		}

		validatorAccountInfo, err := sh.GetSystemAccountInfo(validatorNode)
		if err != nil {
			return err
		}

		validatorLockedBalance, err := sh.getValidatorLockedBalance(validatorNode)
		if err != nil {
			return err
		}

		validatorIdentity, err := sh.getValidatorIdentity(parentIdentity.Account)
		if err != nil {
			return err
		}
		parentName := string(parentIdentity.Name)
		validatorIdentity.DisplayParent = parentName
		validatorIdentity.Parent = parentAccount

		eraStakers, err := sh.getEraStakers(validatorNode)
		if err != nil {
			return err
		}

		nominators := []harvester.Nominator{{
			Address: node,
			Stake:   fmt.Sprint(eraStakers.Own.Int64()),
		}}

		for _, nominee := range eraStakers.Others {
			address, _ := AccountIdToString(nominee.Who, sh.cfg.Name)
			nominators = append(nominators, harvester.Nominator{
				Address: address,
				Stake:   fmt.Sprint(nominee.Value.Int64()),
			})
		}
		validator := harvester.Validator{
			AccountID: node,
			Name:      validatorIdentity.Display,
			Parent: harvester.Parent{
				Name:      parentName,
				AccountID: parentAccount,
			},
			Children:   children,
			Commission: commission,
			Status:     getValidatorStatus(activeValidators, validatorNode),
			Balance:    fmt.Sprint(validatorAccountInfo.Data.Free.Int64()),
			Reserved:   fmt.Sprint(validatorAccountInfo.Data.Reserved.Int64()),
			Locked:     validatorLockedBalance,
			OwnStake:   fmt.Sprint(eraStakers.Own.Int64()),
			TotalStake: fmt.Sprint(eraStakers.Total.Int64()),
			Identity:   validatorIdentity,
			Nominators: nominators,
		}

		// create md5
		validatorJson, err := json.Marshal(validator)
		if err != nil {
			return err
		}

		hashBytes := md5.Sum(validatorJson)
		hash := hex.EncodeToString(hashBytes[:])

		// compare to saved validator map
		dbHash := dbValidatorMap[validator.AccountID]
		if dbHash != hash {
			validator.Hash = hash
			validator.Chain = sh.cfg.Name

			log.Debugf("%s - Saving validator update for account: %s", sh.cfg.Name, node)
			err = sh.repo.SaveValidator(validator)
			if err != nil {
				return err
			}

			log.Debugf("%s - Publishing validator update for account: %s", sh.cfg.Name, node)
			err = sh.publisher.PublishRetained(fmt.Sprintf("harvester/%s/%s/%s", sh.cfg.Name, topic, node), string(validatorJson))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sh *SubstrateHarvester) getCurrentSessionValidators() ([]types.AccountID, error) {
	var validatorAccountIDs []types.AccountID

	key, err := sh.GetStorageDataKey("Session", "Validators", nil)
	if err != nil {
		return validatorAccountIDs, err
	}

	err = sh.GetStorageLatest(key, &validatorAccountIDs)
	if err != nil {
		return nil, err
	}

	return validatorAccountIDs, nil
}

func (sh *SubstrateHarvester) getValidatorPreferences(accountID types.AccountID) (harvester.ValidatorPreferences, error) {
	var validatorPrefs harvester.ValidatorPreferences
	key, err := sh.GetStorageDataKey("Staking", "Validators", accountID[:])
	if err != nil {
		return validatorPrefs, err
	}

	err = sh.GetStorageLatest(key, &validatorPrefs)
	if err != nil {
		return validatorPrefs, err
	}

	return validatorPrefs, nil
}

func (sh *SubstrateHarvester) getValidatorIdentity(accountID types.AccountID) (harvester.ValidatorInfo, error) {
	var identity harvester.ValidatorInfo
	key, err := sh.GetStorageDataKey("Identity", "IdentityOf", accountID[:])
	if err != nil {
		return identity, err
	}

	var validatorIdentity harvester.Identity
	err = sh.GetStorageLatest(key, &validatorIdentity)
	if err != nil {
		return identity, err
	}

	info := harvester.IdentityByte(validatorIdentity.Info.Rest)
	identity = harvester.ValidatorInfo{
		Display:        info.Decode(),
		Legal:          info.Decode(),
		Web:            info.Decode(),
		Riot:           info.Decode(),
		Email:          info.Decode(),
		PgpFingerprint: info.Decode(),
		Image:          info.Decode(),
		Twitter:        info.Decode(),
		Judgements:     validatorIdentity.Judgements,
	}

	return identity, nil
}

func (sh *SubstrateHarvester) getSuperOf(accountID types.AccountID) (harvester.SuperIdentity, error) {
	var superIdentity harvester.SuperIdentity

	key, err := sh.GetStorageDataKey("Identity", "SuperOf", accountID[:])
	if err != nil {
		return superIdentity, err
	}

	err = sh.GetStorageLatest(key, &superIdentity)
	if err != nil {
		return superIdentity, err
	}

	if len(superIdentity.Name) > 0 {
		superIdentity.Name = superIdentity.Name[1:]
	}

	return superIdentity, nil
}

func (sh *SubstrateHarvester) getSubsOf(accountID types.AccountID) (harvester.SubsIdentity, error) {
	var subs harvester.SubsIdentity

	key, err := sh.GetStorageDataKey("Identity", "SubsOf", accountID[:])
	if err != nil {
		return subs, err
	}

	err = sh.GetStorageLatest(key, &subs)
	if err != nil {
		return subs, err
	}

	return subs, nil
}

// get stash accounts
func (sh *SubstrateHarvester) getStashAccounts() ([]string, error) {
	var stashStorageKeys []string

	// Create a key for Staking_Validators (gsrpc does not allow w/o argument, hence override func)
	key, err := sh.CreateStorageKeyUnsafe("Staking", "Validators")
	if err != nil {
		return stashStorageKeys, err
	}

	// This func takes ~4-5 seconds to execute
	stashKeys, err := sh.GetKeysLatest(key)
	if err != nil {
		return stashStorageKeys, err
	}

	for _, key := range stashKeys {
		validatorAddress, err := sh.decodeStashKey(key)
		if err != nil {
			continue // FIXME improve error handling
		}

		stashStorageKeys = append(stashStorageKeys, validatorAddress)
	}

	return stashStorageKeys, nil
}

func (sh *SubstrateHarvester) decodeStashKey(key types.StorageKey) (string, error) {
	keyString := key.Hex()[82:]
	keyBytes, err := types.HexDecodeString(keyString)
	if err != nil {
		return "", err
	}

	validatorAddress, err := AccountIdToString(types.NewAccountID(keyBytes), sh.cfg.Name)
	if err != nil {
		return "", err
	}

	return validatorAddress, nil
}

func getValidatorStatus(parent []types.AccountID, child types.AccountID) string {
	check := Contains(parent, child)
	if check {
		return "active"
	}

	return "candidate"
}

func (sh *SubstrateHarvester) getEraIndex() (types.U32, error) {
	var eraIndex types.U32
	key, err := sh.GetStorageDataKey("Staking", "CurrentEra")
	if err != nil {
		return eraIndex, err
	}

	err = sh.GetStorageLatest(key, &eraIndex)
	if err != nil {
		return eraIndex, err
	}

	return eraIndex, nil
}

func (sh *SubstrateHarvester) getEraStakers(accountID types.AccountID) (types.Exposure, error) {
	var erasStakers types.Exposure

	index, err := sh.getEraIndex()
	if err != nil {
		return erasStakers, err
	}

	eraIndex := UintsToBytes([]uint32{uint32(index)})
	key, err := sh.GetStorageDataKey("Staking", "ErasStakers", eraIndex, accountID[:])
	if err != nil {
		return erasStakers, err
	}

	err = sh.GetStorageLatest(key, &erasStakers)
	if err != nil {
		return erasStakers, err
	}

	return erasStakers, nil
}

func (sh *SubstrateHarvester) getValidatorLockedBalance(accountID types.AccountID) ([]harvester.ValidatorBalancesLocked, error) {
	var balances []harvester.BalancesLocked
	var result []harvester.ValidatorBalancesLocked

	key, err := sh.GetStorageDataKey("Balances", "Locks", accountID[:])
	if err != nil {
		return result, err
	}

	err = sh.GetStorageLatest(key, &balances)
	if err != nil {
		return result, err
	}

	for _, value := range balances {
		result = append(result, harvester.ValidatorBalancesLocked{
			ID:      string(value.ID[:]),
			Amount:  value.Amount.Int64(),
			Reasons: value.Reasons.String(),
		})
	}

	return result, nil
}
