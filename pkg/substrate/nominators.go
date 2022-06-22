package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type PalletStakingNominations struct {
	Targets     []types.AccountID
	submittedIn types.U32
	suppressed  bool
}

func (sh *SubstrateHarvester) getNominatorsCount() (types.U32, error) {
	key, err := sh.GetStorageDataKey("Staking", "CounterForNominators")
	if err != nil {
		return 0, err
	}

	var nominatorsCount types.U32
	err = sh.GetStorageLatest(key, &nominatorsCount)
	if err != nil {
		return 0, err
	}
	return nominatorsCount, nil
}

func (sh *SubstrateHarvester) getNominators() ([]string, error) {
	var stashStorageKeys []string
	key, err := sh.CreateStorageKeyUnsafe("Staking", "Nominators")
	if err != nil {
		return nil, err
	}

	nominatorKeys, err := sh.GetKeysLatest(key)
	if err != nil {
		return nil, err
	}

	for _, key := range nominatorKeys {
		nominatorAddress, err := decodeStashKey(key)
		if err != nil {
			continue // FIXME improve error handling
		}

		stashStorageKeys = append(stashStorageKeys, nominatorAddress)
	}

	return stashStorageKeys, nil
}
