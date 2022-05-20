package repository

import (
	"context"
	"strings"

	"github.com/OdysseyMomentumExperience/harvester/ent/block"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

func (b *Repository) GetBlock(number uint64) (harvester.Block, error) {
	return harvester.Block{}, nil
}

func (b *Repository) SaveBlock(block harvester.Block) error {
	return b.saveBlock(block)
}

func (b *Repository) saveBlock(block harvester.Block) error {
	_, err := b.ent.Block.Create().
		SetNumber(block.Number).
		SetExtrinsicsCount(block.ExtrinsicsCount).
		SetExtrinsics(block.Extrinsics).
		SetFinalized(block.Finalized).
		SetAuthorID(block.AuthorID).
		SetChain(block.Chain).
		Save(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (b *Repository) UpdateFinalizedBlock(bl harvester.Block) error {
	resBlock, err := b.ent.Block.Query().Where(
		block.And(block.ChainEQ(bl.Chain), block.NumberEQ(bl.Number), block.AuthorIDEQ(bl.AuthorID))).Only(context.Background())

	// catch any errors not related to block not found
	if err != nil && !strings.Contains(err.Error(), "block not found") {
		return err
	}

	if resBlock != nil {
		// if block exists update
		_, err = resBlock.Update().
			SetFinalized(bl.Finalized).
			SetExtrinsics(bl.Extrinsics).
			SetExtrinsicsCount(bl.ExtrinsicsCount).
			Save(context.Background())

		if err != nil {
			return err
		}
	} else {
		// else store new block with finalized status
		return b.saveBlock(bl)
	}

	return nil
}
