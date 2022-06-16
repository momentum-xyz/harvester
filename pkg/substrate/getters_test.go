package substrate

import (
	"strconv"
	"strings"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func hexToInt(hex string) (int64, error) {
	s := strings.Replace(hex, "0x", "", -1)
	s = strings.Replace(s, "0X", "", -1)
	return strconv.ParseInt(s, 16, 32)
}

func TestGetters(t *testing.T) {
	t.Run("GetStorageDataKey()", func(t *testing.T) {
		key, err := mockSh.GetStorageDataKey("Session", "CurrentIndex")
		assert.Greater(t, len(key), 0)
		assert.Nil(t, err)

		key, err = mockSh.GetStorageDataKey("a", "b")
		assert.Equal(t, len(key), 0)
		assert.Contains(t, err.Error(), "module a not found")
	})

	t.Run("GetStorageLatest()", func(t *testing.T) {
		key, _ := mockSh.GetStorageDataKey("Staking", "HistoryDepth")
		var historyDepth types.U32
		err := mockSh.GetStorageLatest(key, &historyDepth)
		assert.Greater(t, historyDepth, types.U32(1))
		assert.Nil(t, err)

		var historyDepth2 int64
		err = mockSh.GetStorageLatest(key, &historyDepth2)
		assert.Equal(t, historyDepth2, int64(0))
		assert.NotNil(t, err)

		var historyDepth3 types.U32
		mockSh.GetStorageLatest(nil, &historyDepth3)
		assert.Equal(t, historyDepth3, types.U32(0))
	})

	t.Run("GetChildKeysLatest()", func(t *testing.T) {
		keys, err := mockSh.GetChildKeysLatest(nil, nil)
		assert.Contains(t, err.Error(), "Method not found")
		assert.Equal(t, len(keys), 0)
	})

	t.Run("GetKeysLatest()", func(t *testing.T) {
		key, _ := mockSh.GetStorageDataKey("Staking", "HistoryDepth")
		keys, err := mockSh.GetKeysLatest(key)
		assert.Nil(t, err)
		assert.Greater(t, len(keys), 0)
		assert.IsType(t, &keys[0], &types.StorageKey{})
	})

	t.Run("GetStorageAtBlockHash()", func(t *testing.T) {
		eraDepth := i32tob(0)
		key, _ := mockSh.GetStorageDataKey("System", "BlockHash", eraDepth)
		var blockHash1 types.H256
		err := mockSh.GetStorageLatest(key, &blockHash1)
		blockNumber, _ := hexToInt(blockHash1.Hex())
		assert.Nil(t, err)
		assert.Greater(t, blockNumber, int64(0))

		hash, _ := types.NewHashFromHexString(blockHash1.Hex())
		var blockHash2 types.H256
		err = mockSh.GetStorageAtBlockHash(key, hash, &blockHash2)
		assert.Nil(t, err)
		blockNumber, _ = hexToInt(blockHash2.Hex())
		assert.Nil(t, err)
		assert.Greater(t, blockNumber, int64(0))

		var blockHash3 types.H256
		err = mockSh.GetStorageAtBlockHash(key, types.NewHash([]byte("")), &blockHash3)
		blockNumber, _ = hexToInt(blockHash3.Hex())
		assert.NotNil(t, err)
		assert.Equal(t, blockNumber, int64(0))
	})

	t.Run("QueryStorage()", func(t *testing.T) {
		newHead, _ := mockSh.api.RPC.Chain.GetHeaderLatest()
		finalizedHead, _ := mockSh.api.RPC.Chain.GetFinalizedHead()
		eraDepth := i32tob(0)
		key, _ := mockSh.GetStorageDataKey("System", "BlockHash", eraDepth)
		var blockHash1 types.H256
		err := mockSh.GetStorageLatest(key, &blockHash1)
		blockNumber, _ := hexToInt(blockHash1.Hex())
		assert.Nil(t, err)
		assert.Greater(t, blockNumber, int64(0))

		fromHash, _ := types.NewHashFromHexString(finalizedHead.Hex())
		toHash, _ := types.NewHashFromHexString(finalizedHead.Hex())
		changes, err := mockSh.QueryStorage([]types.StorageKey{key}, fromHash, toHash)
		assert.Nil(t, err)
		assert.IsType(t, changes, []types.StorageChangeSet{})
		assert.Greater(t, len(changes), 0)

		fromHash, _ = types.NewHashFromHexString(newHead.StateRoot.Hex())
		toHash, _ = types.NewHashFromHexString(newHead.StateRoot.Hex())
		changes, err = mockSh.QueryStorage([]types.StorageKey{key}, fromHash, toHash)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), fromHash.Hex())
		assert.Contains(t, err.Error(), toHash.Hex())
		assert.Equal(t, len(changes), 0)
	})

	t.Run("QueryConstant()", func(t *testing.T) {
		var sessionsPerEra types.U32
		err := mockSh.QueryConstant("Staking", "SessionsPerEra", &sessionsPerEra)
		assert.Nil(t, err)
		assert.IsType(t, sessionsPerEra, types.U32(0))

		err = mockSh.QueryConstant("a", "b", &sessionsPerEra)
		assert.NotNil(t, err)
	})

	t.Run("GetSessionLength()", func(t *testing.T) {
		epochDuration, err := mockSh.GetSessionLength()
		assert.Nil(t, err)
		assert.IsType(t, epochDuration, types.U64(0))
	})

	t.Run("GetSlotDuration()", func(t *testing.T) {
		slotDuration, err := mockSh.GetSlotDuration()
		assert.Nil(t, err)
		assert.IsType(t, slotDuration, types.U64(0))
	})

	t.Run("GetSessionPerEra()", func(t *testing.T) {
		sessionsPerEra, err := mockSh.GetSessionPerEra()
		assert.Nil(t, err)
		assert.IsType(t, sessionsPerEra, types.NewU32(0))
	})

	t.Run("GetSessionIndex()", func(t *testing.T) {
		sessionId, err := mockSh.GetSessionIndex()
		assert.Nil(t, err)
		assert.IsType(t, sessionId, types.NewU32(0))
	})

	t.Run("CreateStorageKeyUnsafe()", func(t *testing.T) {
		keys, err := mockSh.CreateStorageKeyUnsafe("Staking", "Validators")
		assert.Nil(t, err)
		assert.Greater(t, len(keys), 1)
	})

	t.Run("getCurrentSessionValidators()", func(t *testing.T) {
		validatorAccountIDs, err := mockSh.getCurrentSessionValidators()
		assert.Nil(t, err)
		assert.IsType(t, &validatorAccountIDs[0], &types.AccountID{})
		assert.Greater(t, len(validatorAccountIDs), 0)
	})
}
