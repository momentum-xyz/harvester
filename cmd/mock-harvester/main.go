package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mock_harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mqtt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/publisher"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"time"
)

const configFileName = "config/mock-harvester/config.dev.yaml"

func main() {
	var err error

	cfg := mock_harvester.GetConfig(configFileName, true)
	mqttClient := mqtt.GetMQTTClient(&cfg.MQTT)

	publisher, err := publisher.NewPublisher(mqttClient)
	if err != nil {
		panic(err)
	}

	for _, chain := range cfg.Chains {
		for topic, attrs := range chain.Topics {
			publishTopic(chain.Name, publisher, topic, attrs)
			if attrs.MsTicker != 0 {
				ticker := time.NewTicker(time.Duration(attrs.MsTicker) * time.Millisecond)
				for {
					select {
					case <-ticker.C:
						publishTopic(chain.Name, publisher, topic, attrs)
					}
				}
			}
		}
	}

}

func publishTopic(chain string, publisher harvester.Publisher, topic string, attr mock_harvester.Topic) {
	value := attr.Value
	if attr.Storage != "" {
		var values []string
		err := json.Unmarshal([]byte(value), &values)
		if err != nil {
			panic(errors.New(fmt.Sprintf("invalid topic value - query storage requires value that can be JSON unmarshalled to []string, original err: %s", err)))
		}
		complementaryValues, err := queryStorage(chain, attr.Storage)
		if err != nil {
			panic(err)
		}
		values = append(values, complementaryValues...)
		valueBytes, err := json.Marshal(values)
		if err != nil {
			panic(err)
		}
		value = string(valueBytes)
	}

	topic = fmt.Sprintf("%s/%s", chain, topic)
	if attr.Retained {
		err := publisher.PublishRetained(topic, value)
		if err != nil {
			panic(err)
		}
	} else {
		err := publisher.Publish(topic, value)
		if err != nil {
			panic(err)
		}
	}
}

func queryStorage(chain string, query string) ([]string, error) {
	rpc, err := getChainRpc(chain)
	if err != nil {
		panic(err)
	}
	api, err := getApi(rpc)
	if err != nil {
		panic(err)
	}
	metadata, err := getMetadata(api)
	if err != nil {
		panic(err)
	}

	switch query {
	case "SocietyMembers":
		return querySocietyMembers(chain, metadata, api)
	default:
		return nil, errors.New(fmt.Sprintf("invalid config - could not match %s with storage query option", query))
	}
}

func querySocietyMembers(chain string, metadata *types.Metadata, api *gsrpc.SubstrateAPI) ([]string, error) {
	key, err := types.CreateStorageKey(metadata, "Society", "Members")
	if err != nil {
		return nil, err
	}

	var target []types.AccountID
	ok, err := api.RPC.State.GetStorageLatest(key, &target)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, err
	}

	networkId := getNetworkID(chain)
	var accountIDs []string
	for _, account := range target {
		accountID, err := subkey.SS58Address(account[:], networkId)
		if err != nil {
			return nil, err
		}
		accountIDs = append(accountIDs, accountID)
	}
	return accountIDs, nil
}

func getChainRpc(chain string) (string, error) {
	switch chain {
	case "kusama":
		return "wss://kusama-rpc.polkadot.io", nil
	case "polkadot":
		return "wss://rpc.polkadot.io", nil
	default:
		return "", errors.New(fmt.Sprintf("invalid configuration - chain name not recognized, incompatible with storage call"))
	}
}

func getApi(rpcAddress string) (*gsrpc.SubstrateAPI, error) {
	return gsrpc.NewSubstrateAPI(rpcAddress)
}

func getMetadata(api *gsrpc.SubstrateAPI) (*types.Metadata, error) {
	return api.RPC.State.GetMetadataLatest()
}

func getNetworkID(name string) uint8 {
	switch name {
	case "kusama":
		return uint8(2)
	case "polkadot":
		return uint8(1)
	default:
		log.Logln(0, "configured network name not recognized, defaulting to kusama")
		return uint8(2)
	}
}
