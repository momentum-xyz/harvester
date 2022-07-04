package substrate

import (
	"testing"

	"github.com/OdysseyMomentumExperience/harvester/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestInflation(t *testing.T) {
	t.Run("GetInflationParams", func(t *testing.T) {
		params, err := mockSh.GetInflationParams(constants.KusamaGenesis)
		assert.Nil(t, err)
		assert.IsType(t, params, InflationParams{})
		assert.Equal(t, params, KusamaInflationParams)

		params, err = mockSh.GetInflationParams(constants.PolkadotGenesis)
		assert.Nil(t, err)
		assert.IsType(t, params, InflationParams{})
		assert.Equal(t, params, PolkadotInflationParams)

		params, err = mockSh.GetInflationParams("")
		assert.Nil(t, err)
		assert.IsType(t, params, InflationParams{})
		assert.Equal(t, params, DefaultInflationParams)
	})

	t.Run("GetInflation", func(t *testing.T) {
		inflation, err := mockSh.GetInflation(1)
		assert.Nil(t, err)
		assert.IsType(t, inflation, Inflation{})
		assert.Greater(t, inflation.StakedFraction, float32(0))
		assert.Greater(t, inflation.StakedReturn, float32(0))

		inflation, err = mockSh.GetInflation(0)
		assert.Nil(t, err)
		assert.IsType(t, inflation, Inflation{})
		assert.Equal(t, inflation.StakedFraction, float32(0))
		assert.Equal(t, inflation.StakedReturn, float32(0))

		inflation, err = mockSh.GetInflation(1000)
		assert.Nil(t, err)
		assert.IsType(t, inflation, Inflation{})
		assert.Greater(t, inflation.StakedFraction, float32(0))
		assert.Greater(t, inflation.StakedReturn, float32(0))
	})
}
