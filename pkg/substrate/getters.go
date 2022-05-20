package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
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

func (sh *SubstrateHarvester) GetStorageAtBlockHash(key types.StorageKey, hash types.Hash, target interface{}) error {
	ok, err := sh.api.RPC.State.GetStorage(key, target, hash)
	if err != nil {
		return (err)
	} else if !ok {
		return err
	}

	return nil
}

func (sh *SubstrateHarvester) getCurrentSessionValidators() ([]types.AccountID, error) {
	key, err := sh.GetStorageDataKey("Session", "Validators", nil)
	if err != nil {
		return nil, err
	}

	var validatorAccountIDs []types.AccountID
	err = sh.GetStorageLatest(key, &validatorAccountIDs)
	if err != nil {
		return nil, err
	}

	return validatorAccountIDs, nil
}

// type CustomSignedBlock struct {
// 	*types.SignedBlock
// }

// type ID types.AccountID
// type Node struct {
// 	AccountID ID
// 	Index     int
// }

// type CustomBlock struct {
// 	Number    uint32
// 	AuthorID  string
// 	Finalized bool
// 	Extrinsic types.Extrinsic
// }

// var AuthIndexes = make(map[string]Node)
// var StringIDMap = make(map[string]ID)

// func (block *CustomSignedBlock) ToFormat() CustomBlock {
// 	return CustomBlock{
// 		Number: uint32(block.Block.Header.Number),
// 	}
// }

// func (validator ID) ToString() string {
// 	id, err := subkey.SS58Address(validator[:], 2)
// 	if err != nil {
// 		panic(err)
// 	}

// 	StringIDMap[id] = validator
// 	return id
// }

// func (validator ID) ToAccount() types.AccountInfo {
// 	var target types.AccountInfo
// 	GetStorageData("System", "Account", &target, validator[:])

// 	return target
// }

// func (validator ID) ToBonded() string {
// 	var target ID
// 	GetStorageData("Staking", "Bonded", &target, validator[:])

// 	return target.ToString()
// }

// func GetMetadata() (*gsrpc.SubstrateAPI, *types.Metadata) {
// 	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	metadata, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return api, metadata
// }

// 	ok, err := api.RPC.State.GetStorageLatest(key, target)
// 	if err != nil {
// 		panic(err)
// 	} else if !ok {
// 		return nil
// 	}

// 	return target
// }

//func GetValidators() []string {
//	var target []ID
//	var result []string
//	GetStorageData("Session", "Validators", &target)
//
//	if len(target) > 0 {
//		result = make([]string, len(target))
//
//		for index, v := range target {
//			hex := v.ToString()
//			result[index] = hex
//			AuthIndexes[hex] = Node{v, index}
//		}
//	}
//
//	return result
//}

// func GetBabeAuthorities() []string {
// 	var target []ID
// 	var result []string
// 	GetStorageData("Babe", "Authorities", &target)

// 	if len(target) > 0 {
// 		result = make([]string, len(target))

// 		for index, v := range target {
// 			hex := v.ToString()
// 			result[index] = hex
// 			AuthIndexes[hex] = Node{v, index}
// 		}
// 	}

// 	return result
// }

// func GetBabeNextAuthorities() []string {
// 	var target []ID
// 	var result []string
// 	GetStorageData("Babe", "NextAuthorities", &target)

// 	if len(target) > 0 {
// 		result = make([]string, len(target))

// 		for index, v := range target {
// 			result[index] = v.ToString()
// 		}
// 	}

// 	return result
// }

// func GetLatestAuthor() string {
// 	var target ID
// 	GetStorageData("Authorship", "Author", &target)

// 	return target.ToString()
// }

// func GetSessionID() types.U32 {
// 	var target types.U32
// 	GetStorageData("Session", "CurrentIndex", &target)

// 	return target
// }

// func GetAuthoredBlocks(sessionId types.U32, validator types.AccountID) types.U32 {
// 	var target types.U32
// 	byteId, err := sessionId.MarshalJSON()
// 	fmt.Println(byteId)
// 	if err != nil {
// 		panic(err)
// 	}

// 	GetStorageData("ImOnline", "AuthoredBlocks", &target, byteId, validator[:])

// 	return target
// }

// func GetReceivedHeartbeats(sessionId types.U32, validator ID) types.U32 {
// 	var target types.U32
// 	byteId := byte(sessionId)
// 	fmt.Println(byteId)
// 	validatorId := byte(AuthIndexes[validator.ToString()].Index)
// 	fmt.Println(validatorId)

// 	GetStorageData("ImOnline", "ReceivedHeartbeats", &target, []byte{byteId}, []byte{validatorId})

// 	return target
// }

// 	return CustomSignedBlock(latest_block)
// }
