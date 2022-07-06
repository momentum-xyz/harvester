package substrate

import (
	"reflect"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestNominators(t *testing.T) {
	t.Run("getActiveNominators()", func(t *testing.T) {
		nominators, err := mockSh.getActiveNominators(func(err error) {})
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(nominators), reflect.SliceOf(reflect.TypeOf(types.NewAccountID([]byte{}))))
	})
}
