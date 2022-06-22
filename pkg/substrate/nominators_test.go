package substrate

import (
	"reflect"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestNominators(t *testing.T) {
	t.Run("getNominatorsCount()", func(t *testing.T) {
		nominatorsCount, err := sh.getNominatorsCount()
		assert.Nil(t, err)
		assert.IsType(t, nominatorsCount, types.NewU32(0))
	})

	t.Run("getNominators()", func(t *testing.T) {
		nominators, err := sh.getNominators()
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(nominators), reflect.SliceOf(reflect.TypeOf("")))
	})
}
