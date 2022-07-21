package substrate

import (
	"math"

	"github.com/OdysseyMomentumExperience/harvester/pkg/constants"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/xxhash"
	"github.com/pkg/errors"
)

type SystemProperties struct {
	SS58Format    uint8
	TokenDecimals uint32
	TokenSymbol   string
}

func (sh *SubstrateHarvester) GetStorageDataKey(prefix string, method string, args ...[]byte) (types.StorageKey, error) {
	key, err := types.CreateStorageKey(sh.metadata, prefix, method, args...)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (sh *SubstrateHarvester) GetStorageLatest(key types.StorageKey, target interface{}) error {
	ok, err := sh.api.RPC.State.GetStorageLatest(key, target)
	if err != nil {
		return err
	} else if !ok {
		return err
	}

	return nil
}

func (sh *SubstrateHarvester) GetChildKeysLatest(childKey, prefix types.StorageKey) ([]types.StorageKey, error) {
	keys, err := sh.api.RPC.State.GetChildKeysLatest(childKey, prefix)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (sh *SubstrateHarvester) GetKeysLatest(key types.StorageKey) ([]types.StorageKey, error) {
	keys, err := sh.api.RPC.State.GetKeysLatest(key)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (sh *SubstrateHarvester) GetStorageAtBlockHash(key types.StorageKey, hash types.Hash, target interface{}) error {
	ok, err := sh.api.RPC.State.GetStorage(key, target, hash)
	if err != nil || !ok {
		return err
	}

	return nil
}

func (sh *SubstrateHarvester) QueryStorage(keys []types.StorageKey, fromBlock types.Hash, toBlock types.Hash) ([]types.StorageChangeSet, error) {
	changes, err := sh.api.RPC.State.QueryStorage(keys, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}

	return changes, nil
}

func (sh *SubstrateHarvester) QueryStorageAt(keys []types.StorageKey) ([]types.KeyValueOption, error) {
	hexKeys := make([]string, len(keys))
	for i, key := range keys {
		hexKeys[i] = key.Hex()
	}

	var res []types.StorageChangeSet
	err := sh.api.Client.Call(&res, "state_queryStorageAt", hexKeys)
	if err != nil {
		return nil, err
	}

	var changes []types.KeyValueOption
	for _, r := range res {
		changes = append(changes, r.Changes...)
	}
	return changes, nil
}

func (sh *SubstrateHarvester) QueryConstant(prefix string, name string, res interface{}) error {
	data, err := sh.metadata.FindConstantValue(prefix, name)
	if err != nil {
		return err
	}
	return types.DecodeFromBytes(data, res)
}

func (sh *SubstrateHarvester) GetSessionLength() (types.U64, error) {

	var epochDuration types.U64

	err := sh.QueryConstant("Babe", "EpochDuration", &epochDuration)

	if err != nil {
		return 0, errors.Wrap(err, "error while fetching epoch duration constant")
	}
	return epochDuration, nil
}

func (sh *SubstrateHarvester) GetSlotDuration() (types.U64, error) {
	var slotDuration types.U64

	err := sh.QueryConstant("Babe", "ExpectedBlockTime", &slotDuration)

	if err != nil {
		return 0, errors.Wrap(err, "error while fetching slot duration constant")
	}
	return slotDuration, nil
}

func (sh *SubstrateHarvester) GetSessionPerEra() (types.U32, error) {
	var sessionsPerEra types.U32

	err := sh.QueryConstant("Staking", "SessionsPerEra", &sessionsPerEra)

	if err != nil {
		return 0, errors.Wrap(err, "error while fetching session per era constant value")
	}
	return sessionsPerEra, nil
}

func (sh *SubstrateHarvester) GetSessionIndex() (types.U32, error) {
	var sessionId types.U32

	key, err := sh.GetStorageDataKey("Session", "CurrentIndex")
	if err != nil {
		return sessionId, err
	}

	err = sh.GetStorageLatest(key, &sessionId)
	if err != nil {
		return sessionId, err
	}

	return sessionId, nil
}

func (sh *SubstrateHarvester) CreateStorageKeyUnsafe(prefix string, method string, args ...[]byte) (types.StorageKey, error) {
	key := append(xxhash.New128([]byte(prefix)).Sum(nil), xxhash.New128([]byte(method)).Sum(nil)...)

	if len(args) > 0 {
		entryMeta, err := sh.FindStorageEntryMetadata(prefix, method)
		if err != nil {
			return nil, err
		}
		hashers, err := entryMeta.Hashers()
		if err != nil {
			return nil, err
		}

		for i, arg := range args {
			_, err := hashers[i].Write(arg)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to hash args[%d]: %s Error: %v", i, arg, err)
			}
			key = append(key, hashers[i].Sum(nil)...)
		}

	}
	return key, nil
}

func (sh *SubstrateHarvester) FindStorageEntryMetadata(prefix string, method string) (types.StorageEntryMetadata, error) {
	meta, err := sh.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	entryMeta, err := meta.FindStorageEntryMetadata(prefix, method)
	if err != nil {
		return nil, err
	}
	return entryMeta, nil
}

func (sh *SubstrateHarvester) GetTotalIssuance() (types.U128, error) {
	var totalIssuance types.U128
	key, err := sh.GetStorageDataKey("Balances", "TotalIssuance")
	if err != nil {
		return totalIssuance, err
	}

	err = sh.GetStorageLatest(key, &totalIssuance)
	if err != nil {
		return totalIssuance, err
	}

	return totalIssuance, nil
}

func (sh *SubstrateHarvester) GetAuctionCounter() (types.U32, error) {
	var auctionCounter types.U32
	key, err := sh.GetStorageDataKey("Auctions", "AuctionCounter")
	if err != nil {
		return auctionCounter, err
	}

	err = sh.GetStorageLatest(key, &auctionCounter)
	if err != nil {
		return auctionCounter, err
	}

	return auctionCounter, nil
}

func (sh *SubstrateHarvester) GetGenesisHash() (types.Hash, error) {
	hash, err := sh.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

func (sh *SubstrateHarvester) GetSystemProperties() (SystemProperties, error) {
	var p SystemProperties
	err := sh.api.Client.Call(&p, "system_properties")
	return p, err
}

func (sh *SubstrateHarvester) GetSystemAccountInfo(accountID types.AccountID) (harvester.AccountInfo, error) {
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

func (sh *SubstrateHarvester) GetAccountBalance(accountID types.AccountID) (float64, error) {
	accountInfo, err := sh.GetSystemAccountInfo(accountID)
	if err != nil || accountInfo.Data.Free.Int == nil {
		return 0, err
	}

	systemProperties, err := sh.GetSystemProperties()
	if err != nil {
		return 0, err
	}

	decimals := systemProperties.TokenDecimals
	if decimals == 0 {
		decimals = constants.DefaultTokenDecimals
	}

	balance := float64(accountInfo.Data.Free.Int64()) / math.Pow10(int(decimals))

	return balance, nil
}
