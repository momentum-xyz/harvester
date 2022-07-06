package substrate

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/decred/base58"
	"github.com/vedhavyas/go-subkey"
)

func Contains[T comparable](parent []T, child T) bool {
	for _, value := range parent {
		if value == child {
			return true
		}
	}

	return false
}

func UintsToBytes(vs []uint32) []byte {
	buf := make([]byte, len(vs)*4)
	for i, v := range vs {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}

func AccountIdToString(id types.AccountID, network ...string) (string, error) {
	networkId := uint8(2) // default network
	if len(network) > 0 {
		networkId = getNetworkID(network[0])
	}
	address, err := subkey.SS58Address(id[:], networkId)
	if err != nil {
		return "", err
	}

	return address, nil
}

func StringToAccountId(account string) (types.AccountID, error) {
	addressBytes := base58.Decode(account)
	publicKey := addressBytes[1 : len(addressBytes)-2]
	len := len(publicKey)
	if len != 32 {
		return types.NewAccountID(nil), fmt.Errorf("%s address yielded wrong length", account)
	}

	return types.NewAccountID(publicKey), nil
}

func Round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
