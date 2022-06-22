package substrate

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/log"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/decred/base58"
	"github.com/vedhavyas/go-subkey"
)

func (sh *SubstrateHarvester) ProcessValidators(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {
	log.Debug("processing validators")

	ticker := time.NewTicker(5 * time.Minute)
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
		validatorNode, err := stringToAccountId(node)
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

		parentAccount, err := accountIdToString(parentIdentity.Account)
		if err != nil {
			return err
		}

		subAccounts, err := sh.getSubsOf(validatorNode)
		if err != nil {
			return err
		}

		children := make([]string, 0)
		for _, c := range subAccounts.Accounts {
			childAccount, err := accountIdToString(c)
			if err != nil {
				return err
			}
			children = append(children, childAccount)
		}

		validatorAccountInfo, err := sh.getSystemAccountInfo(validatorNode)
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
			address, _ := accountIdToString(nominee.Who)
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
		validatorAddress, err := decodeStashKey(key)
		if err != nil {
			continue // FIXME improve error handling
		}

		stashStorageKeys = append(stashStorageKeys, validatorAddress)
	}

	return stashStorageKeys, nil
}

func accountIdToString(id types.AccountID) (string, error) {
	address, err := subkey.SS58Address(id[:], 2)
	if err != nil {
		return "", err
	}

	return address, nil
}

func stringToAccountId(account string) (types.AccountID, error) {
	addressBytes := base58.Decode(account)
	publicKey := addressBytes[1 : len(addressBytes)-2]
	len := len(publicKey)
	if len != 32 {
		return types.NewAccountID(nil), fmt.Errorf("%s address yielded wrong length", account)
	}

	return types.NewAccountID(publicKey), nil
}

func decodeStashKey(key types.StorageKey) (string, error) {
	keyString := key.Hex()[82:]
	keyBytes, err := types.HexDecodeString(keyString)
	if err != nil {
		return "", err
	}

	validatorAddress, err := accountIdToString(types.NewAccountID(keyBytes))
	if err != nil {
		return "", err
	}

	return validatorAddress, nil
}

func getValidatorStatus(parent []types.AccountID, child types.AccountID) string {
	check := contains(parent, child)
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

	eraIndex := uintsToBytes([]uint32{uint32(index)})
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

func (sh *SubstrateHarvester) getSystemAccountInfo(accountID types.AccountID) (harvester.AccountInfo, error) {
	var accountInfo harvester.AccountInfo

	key, err := sh.GetStorageDataKey("System", "Account", accountID[:])
	if err != nil {
		return accountInfo, err
	}

	err = sh.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return accountInfo, err
	}

	return accountInfo, nil
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

func contains(parent []types.AccountID, child types.AccountID) bool {
	for _, value := range parent {
		if value == child {
			return true
		}
	}

	return false
}

func uintsToBytes(vs []uint32) []byte {
	buf := make([]byte, len(vs)*4)
	for i, v := range vs {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}
