package substrate

import (
	"encoding/binary"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("Contains()", func(t *testing.T) {
		accountsId := []types.AccountID{types.NewAccountID([]byte{1}), types.NewAccountID([]byte{2})}
		assert.True(t, Contains(accountsId, types.NewAccountID([]byte{1})))
		assert.False(t, Contains(accountsId, types.NewAccountID([]byte{100})))
	})

	t.Run("UintsToBytes()", func(t *testing.T) {
		a, b := make([]byte, 4), make([]byte, 4)
		binary.LittleEndian.PutUint32(a, 1)
		binary.LittleEndian.PutUint32(b, 2)
		assert.Equal(t, UintsToBytes([]uint32{1, 2}), append(a, b...))
		assert.NotEqual(t, UintsToBytes([]uint32{1, 2, 3}), append(a, b...))
	})

	t.Run("AccountIdToString()", func(t *testing.T) {
		a := types.NewAddressFromAccountID([]byte{1})
		accountIdStr, err := AccountIdToString(a.AsAccountID, mockSh.cfg.Name)
		accountId, _ := StringToAccountId(accountIdStr)
		assert.Equal(t, a.AsAccountID, accountId)
		assert.Nil(t, err)
	})

	t.Run("StringToAccountId()", func(t *testing.T) {
		a := types.NewAddressFromAccountID([]byte{2})
		accountIdStr, _ := AccountIdToString(a.AsAccountID, mockSh.cfg.Name)
		accountId, err := StringToAccountId(accountIdStr)
		assert.Equal(t, a.AsAccountID, accountId)
		assert.Nil(t, err)

		accountId, err = StringToAccountId("xxx")
		assert.Equal(t, types.NewAccountID(nil), accountId)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "xxx address yielded wrong length")
	})

	t.Run("Round()", func(t *testing.T) {
		assert.Equal(t, Round(12.578, 2), 12.58)
		assert.NotEqual(t, Round(12.578, 3), 12.58)
		assert.Equal(t, Round(12.573, 2), 12.57)
		assert.NotEqual(t, Round(12.573, 2), 12.58)
	})
}
