package substrate

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestChainEraState(t *testing.T) {
	t.Run("GetChainEraState", func(t *testing.T) {
		activeEra, _ := sh.GetActiveEra()
		err := sh.GetChainEraState(activeEra, uint32(1), uint32(1), func(err error) {})
		assert.Nil(t, err)
	})

	t.Run("GetErasValidatorReward", func(t *testing.T) {
		activeEra, _ := sh.GetActiveEra()
		erasValidatorReward, err := sh.GetErasValidatorReward(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(erasValidatorReward), reflect.TypeOf(types.NewU128(*big.NewInt(0))))
	})

	t.Run("getEraTotalStake", func(t *testing.T) {
		activeEra, _ := sh.GetActiveEra()
		totalStakeInEra, err := sh.getEraTotalStake(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(totalStakeInEra), reflect.TypeOf(types.NewU128(*big.NewInt(0))))
	})

}
