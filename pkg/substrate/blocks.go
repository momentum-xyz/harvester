package substrate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey"
)

type RawBabePreDigestCompat struct {
	IsZero  bool
	AsZero  types.U32
	IsOne   bool
	AsOne   types.U32
	IsTwo   bool
	AsTwo   types.U32
	IsThree bool
	AsThree types.U32
}

func (sh *SubstrateHarvester) ProcessNewHeader(ctx context.Context, fn harvester.ErrorHandler) error {

	log.Debug("proccesing new headers")

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	var latestNewHead types.U32 = 0

	for {
		var err error
		select {
		case <-ticker.C:
			latestNewHead, err = sh.fetchNewHead(latestNewHead)
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) ProcessFinalizedHeader(ctx context.Context, fn harvester.ErrorHandler) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var latestFinalizedHead types.U32 = 0

	for {
		var err error
		select {
		case <-ticker.C:
			latestFinalizedHead, err = sh.fetchFinalizedHead(latestFinalizedHead)

		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}
}

func (sh *SubstrateHarvester) fetchNewHead(latestHead types.U32) (types.U32, error) {

	var nextNewdHead types.U32

	newHead, err := sh.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return latestHead, errors.Wrapf(err, " error while fetching latest header, previous block is %v", latestHead)
	}

	nextNewdHead = types.U32(newHead.Number)
	if nextNewdHead == latestHead {
		log.Debug("no new header, latest is :", latestHead)
		return latestHead, nil
	}
	err = sh.processHeader(newHead, "block-creation-event", false)
	if err != nil {
		return latestHead, errors.Wrapf(err, " error while processing new head %v", nextNewdHead)
	}

	return nextNewdHead, nil
}

func (sh *SubstrateHarvester) fetchFinalizedHead(latestHead types.U32) (types.U32, error) {

	var nextFinalizedHead types.U32

	finalizedHead, err := sh.api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return latestHead, errors.Wrapf(err, " error while fetching finalizedHead hash for block  %v", latestHead)
	}
	finalizedHeader, err := sh.api.RPC.Chain.GetHeader(finalizedHead)
	if err != nil {
		return latestHead, errors.Wrapf(err, " error while fetching head for block %v", latestHead)
	}

	nextFinalizedHead = types.U32(finalizedHeader.Number)
	if nextFinalizedHead == latestHead {
		log.Debug("no new finalized block latest is :", latestHead)
		return latestHead, nil
	}
	err = sh.processHeader(finalizedHeader, "block-finalized-event", true)
	if err != nil {
		return latestHead, errors.Wrapf(err, " error while processing finalized head %v", nextFinalizedHead)
	}

	return nextFinalizedHead, nil
}

func (sh *SubstrateHarvester) processHeader(header *types.Header, topic string, finalized bool) error {
	blockNumber := uint64(header.Number)
	hash, err := sh.api.RPC.Chain.GetBlockHash(blockNumber)
	if err != nil {
		return errors.Wrapf(err, "getBlockHash for block number %v is failed", blockNumber)
	}

	signedBlock, err := sh.api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return errors.Wrapf(err, "getBlock for hash %v is failed", spew.Sdump(hash))
	}

	var res types.U32
	for _, di := range header.Digest {
		if di.IsPreRuntime {
			var digest RawBabePreDigestCompat
			//TODO Add di.AsPreRuntime.Bytes in logs
			decoder := scale.NewDecoder(bytes.NewReader(di.AsPreRuntime.Bytes))
			err := decoder.Decode(&digest)
			if err != nil {
				return errors.Wrapf(err, "header digest decoding for block number %v is failed", blockNumber)
			}
			res = digest.AsZero
		}
	}

	validatorAccountIDs, err := sh.getCurrentSessionValidators()
	if err != nil {
		return errors.Wrapf(err, "header digest decoding for block number %v is failed", blockNumber)
	}
	auth := validatorAccountIDs[res]
	networkId := getNetworkID(sh.cfg.Name)
	validatorID, err := subkey.SS58Address(auth[:], networkId)
	if err != nil {
		return errors.Wrapf(err, "SS58Address validatorID for block  number %v is failed", blockNumber)
	}
	extrinsicsCount := len(signedBlock.Block.Extrinsics)
	err = sh.publishHeader(header, extrinsicsCount, validatorID, topic, finalized)
	if err != nil {
		return errors.Wrapf(err, "publish header info for topic %v and validator %v is failed", topic, validatorID)
	}
	return nil
}
func (sh *SubstrateHarvester) publishHeader(
	header *types.Header,
	extrinsicsCount int,
	validatorID string,
	headerTopic string,
	finalized bool) error {
	block := harvester.Block{
		Number:          uint32(header.Number),
		AuthorID:        validatorID,
		Finalized:       finalized,
		ExtrinsicsCount: extrinsicsCount,
		Chain:           sh.cfg.Name,
	}

	blockJson, err := json.Marshal(block)
	if err != nil {
		return err
	}

	log.Debugf("%s - Publishing BlockCreationEvent for block number: # %d", sh.cfg.Name, block.Number)
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/%s", sh.cfg.Name, headerTopic), string(blockJson))
	if err != nil {
		return err
	}

	log.Debugf("%s - Storing Block For %s : # %d", sh.cfg.Name, headerTopic, block.Number)
	err = sh.repo.SaveBlock(block)
	if err != nil {
		return err
	}
	return nil
}
