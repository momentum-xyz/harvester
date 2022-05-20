package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
)

type ID types.AccountID

func (validator *ID) ToString() string {
	id, err := subkey.SS58Address(validator[:], 2)
	if err != nil {
		panic(err)
	}

	return id
}

func (validator *ID) ToAccount() types.AccountInfo {
	var target types.AccountInfo
	//GetStorageData("System", "Account", &target, validator[:])

	return target
}

func (validator *ID) ToBonded() string {
	var target ID
	//GetStorageData("Staking", "Bonded", &target, validator[:])

	return target.ToString()
}

func GetValidators() []ID {
	var target []ID
	//GetStorageData("Session", "Validators", &target)

	return target
}

func GetLatestAuthor() ID {
	var target ID
	//GetStorageData("Authorship", "Author", &target)

	return target
}
