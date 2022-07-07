package substrate

import (
	"math"

	"github.com/OdysseyMomentumExperience/harvester/pkg/constants"
)

type InflationParams struct {
	AuctionAdjust float32
	AuctionMax    float32
	Falloff       float32
	MaxInflation  float32
	MinInflation  float32
	StakeTarget   float32
}

type Inflation struct {
	IdealStake     float32
	IdealInterest  float32
	Inflation      float32
	StakedFraction float32
	StakedReturn   float32
}

var DefaultInflationParams = InflationParams{
	AuctionAdjust: 0,
	AuctionMax:    0,
	Falloff:       0.05,
	MaxInflation:  0.1,
	MinInflation:  0.025,
	StakeTarget:   0.5,
}

var KusamaInflationParams = InflationParams{
	AuctionAdjust: (0.3 / 60),
	AuctionMax:    60,
	Falloff:       DefaultInflationParams.Falloff,
	MaxInflation:  DefaultInflationParams.MaxInflation,
	MinInflation:  DefaultInflationParams.MinInflation,
	StakeTarget:   0.75,
}

var PolkadotInflationParams = InflationParams{
	AuctionAdjust: DefaultInflationParams.AuctionAdjust,
	AuctionMax:    DefaultInflationParams.AuctionMax,
	Falloff:       DefaultInflationParams.Falloff,
	MaxInflation:  DefaultInflationParams.MaxInflation,
	MinInflation:  DefaultInflationParams.MinInflation,
	StakeTarget:   0.75,
}

func (sh *SubstrateHarvester) GetInflationParams(genesisHash string) (InflationParams, error) {
	switch genesisHash {
	case constants.KusamaGenesis:
		return KusamaInflationParams, nil
	case constants.PolkadotGenesis:
		return PolkadotInflationParams, nil
	}
	return DefaultInflationParams, nil
}

func (sh *SubstrateHarvester) GetInflation(stakingRatio float32) (Inflation, error) {
	var inflation Inflation
	genesisHash, err := sh.GetGenesisHash()
	if err != nil {
		return inflation, err
	}

	params, err := sh.GetInflationParams(genesisHash.Hex())
	if err != nil {
		return inflation, err
	}

	numAuctions, err := sh.GetAuctionCounter()
	if err != nil {
		return inflation, err
	}

	inflation.StakedFraction = float32(Round(float64(stakingRatio)/100, 6))
	inflation.IdealStake = func() float32 {
		if params.AuctionMax < float32(numAuctions) {
			return params.StakeTarget - (params.AuctionMax * params.AuctionAdjust)
		}
		return params.StakeTarget - (float32(numAuctions) * params.AuctionAdjust)
	}()
	inflation.IdealInterest = params.MaxInflation / inflation.IdealStake
	inflation.Inflation = 100 * (params.MinInflation + func() float32 {
		if inflation.StakedFraction <= inflation.IdealStake {
			return (inflation.StakedFraction * (inflation.IdealInterest - (params.MinInflation / inflation.IdealStake)))
		}
		return (((inflation.IdealInterest * inflation.IdealStake) - params.MinInflation) * float32(math.Pow(2, (float64(inflation.IdealStake)-float64(inflation.StakedFraction))/float64(params.Falloff))))
	}())
	inflation.StakedReturn = func() float32 {
		if inflation.StakedFraction > 0 {
			return (inflation.Inflation) / (inflation.StakedFraction)
		}
		return 0
	}()

	// convert IdealStake to percentage
	inflation.IdealStake = inflation.IdealStake * 100

	return inflation, err
}
