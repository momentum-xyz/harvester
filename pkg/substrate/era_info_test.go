package substrate

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEraInfo(t *testing.T) {
	t.Run("GetActiveEra", func(t *testing.T) {
		activeEra, err := mockSh.GetActiveEra()
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(activeEra).String(), reflect.Uint32.String())
	})

	t.Run("GetActiveEraDepth", func(t *testing.T) {
		activeEraDepth, err := mockSh.GetActiveEraDepth()
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(activeEraDepth), reflect.TypeOf([]byte("")))
	})

	t.Run("GetEraDepth", func(t *testing.T) {
		activeEra, _ := mockSh.GetActiveEra()
		eraDepth, err := mockSh.GetEraDepth(activeEra)
		assert.Nil(t, err)
		assert.IsType(t, reflect.TypeOf(eraDepth), reflect.TypeOf([]byte("")))
	})

}
