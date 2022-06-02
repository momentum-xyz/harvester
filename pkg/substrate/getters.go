package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/xxhash"
	"github.com/pkg/errors"
)

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
	if err != nil {
		return err
	} else if !ok {
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

func (sh *SubstrateHarvester) CreateStorageKeyUnsafe(prefix string, method string) (types.StorageKey, error) {
	return append(xxhash.New128([]byte(prefix)).Sum(nil), xxhash.New128([]byte(method)).Sum(nil)...), nil
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
