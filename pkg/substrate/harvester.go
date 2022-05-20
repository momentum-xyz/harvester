package substrate

import (
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type SubstrateHarvester struct {
	cfg       harvester.ChainConfig
	publisher harvester.Publisher
	repo      harvester.Repository
	api       *gsrpc.SubstrateAPI
	metadata  *types.Metadata
}

func NewHarvester(cfg harvester.ChainConfig,
	pub harvester.Publisher,
	repo harvester.Repository) (*SubstrateHarvester, error) {
	api, err := getApi(cfg.RPC)
	if err != nil {
		return nil, err
	}

	metadata, err := getMetadata(api)
	if err != nil {
		return nil, err
	}

	return &SubstrateHarvester{
		cfg:       cfg,
		publisher: pub,
		repo:      repo,
		api:       api,
		metadata:  metadata,
	}, nil
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

func (sh *SubstrateHarvester) getActiveProcesses() []harvester.ActiveHarvesterProcess {

	var substrateProcesses []harvester.ActiveHarvesterProcess
	for _, activeProcess := range sh.cfg.ActiveTopics {
		substrateProcesses = append(substrateProcesses, sh.topicProcessorStore()(activeProcess))
	}
	return substrateProcesses
}

func (sh *SubstrateHarvester) topicProcessorStore() func(string) harvester.ActiveHarvesterProcess {

	innerProcessStore := map[string]harvester.ActiveHarvesterProcess{
		"block-creation-event":  sh.ProcessNewHeader,
		"block-finalized-event": sh.ProcessFinalizedHeader,
		"reward-event":          sh.ProcessErasRewardPoints,
		"society-members":       sh.ProcessSocietyMembers,
		"extrinsics-pool":       sh.ProcessPendingExtrinsics,
	}

	return func(topic string) harvester.ActiveHarvesterProcess {
		return innerProcessStore[topic]
	}
}

func (sh *SubstrateHarvester) Stop() {
	log.Errorf("stopping chain harvester %s\n", sh.cfg.Name)
}
