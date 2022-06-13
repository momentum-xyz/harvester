package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

type ExtrinsicItem struct {
	Method   string
	IsSigned bool
}

type PendingExtrinsics struct {
	Count int `json:"pendingExtrinsics"`
}

func (sh *SubstrateHarvester) ProcessPendingExtrinsics(ctx context.Context,
	fn harvester.ErrorHandler,
	pmc harvester.PerformanceMonitorClient,
	topic string) error {

	ticker := time.NewTicker(3 * time.Second)

	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			err = sh.publishPendingExtrinsics(fn, pmc, topic)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}

}

func (sh *SubstrateHarvester) publishPendingExtrinsics(fn harvester.ErrorHandler, pmc harvester.PerformanceMonitorClient, topic string) error {
	defer pmc.WriteProcessResponseMetrics(time.Now(), topic, fn)

	pendingExtrinsics, err := sh.api.RPC.Author.PendingExtrinsics()
	if err != nil {
		return errors.Wrap(err, "error occurred while fetching rpc author pending extrinsics")
	}

	extrinsicsCount := PendingExtrinsics{
		Count: len(pendingExtrinsics),
	}

	msg, err := json.Marshal(extrinsicsCount)
	if err != nil {
		return errors.Wrapf(err, "error occurred on json marshal extrinsicsCount %v", extrinsicsCount)
	}

	log.Debugf("%s - Publishing ExtrinsicsPool update", sh.cfg.Name)
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/%s", sh.cfg.Name, topic), string(msg))
	if err != nil {
		return errors.Wrap(err, "error occurred while publishing extrinsicsCount")
	}
	return nil
}

func (sh *SubstrateHarvester) ProcessExtrinsics(rawExtrinsics []types.Extrinsic) []harvester.Extrinsic {
	var extrinsics []harvester.Extrinsic
	for _, ext := range rawExtrinsics {
		section, method := FindCallValues(sh.metadata, ext.Method.CallIndex)
		extrinsics = append(extrinsics, harvester.Extrinsic{
			Method: section + "." + method,
		})
	}
	return extrinsics
}

func FindCallValues(metadata *types.Metadata, call types.CallIndex) (string, string) {
	var sectionName string
	var methodName string

	for _, mod := range metadata.AsMetadataV14.Pallets {
		if mod.HasCalls && uint8(mod.Index) == call.SectionIndex {
			sectionName = string(mod.Name)
			callType := mod.Calls.Type.Int64()

			if typ, ok := metadata.AsMetadataV14.EfficientLookup[callType]; ok {
				index := call.MethodIndex
				methodName = string(typ.Def.Variant.Variants[index].Name)
			}

			break
		}
	}

	return sectionName, methodName
}
