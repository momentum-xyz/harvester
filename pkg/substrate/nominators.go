package substrate

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func (sh *SubstrateHarvester) getActiveNominators(fn harvester.ErrorHandler) ([]types.AccountID, error) {
	eraDepth, _ := sh.GetActiveEraDepth()
	k, err := sh.CreateStorageKeyUnsafe("Staking", "ErasStakers", eraDepth)
	if err != nil {
		return nil, err
	}

	stashKeys, err := sh.GetKeysLatest(k)
	if err != nil {
		return nil, err
	}

	changes, err := sh.QueryStorageAt(stashKeys)
	if err != nil {
		return nil, err
	}

	var nominators []types.AccountID
	for _, change := range changes {
		var r types.Exposure
		err = types.DecodeFromBytes(change.StorageData, &r)
		if err != nil {
			return nil, err
		}

		for _, n := range r.Others {
			if !Contains(nominators, n.Who) {
				nominators = append(nominators, n.Who)
			}
		}
	}

	return nominators, nil
}
