package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestValidators(t *testing.T) {
	t.Run("getCurrentSessionValidators()", func(t *testing.T) {
		validatorAccountIDs, err := mockSh.getCurrentSessionValidators()
		assert.Nil(t, err)
		assert.IsType(t, &validatorAccountIDs[0], &types.AccountID{})
		assert.Greater(t, len(validatorAccountIDs), 0)
	})
}
