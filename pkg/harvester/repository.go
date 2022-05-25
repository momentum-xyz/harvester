package harvester

type Repository interface {
	GetBlock(number uint64) (Block, error)
	SaveBlock(Block) error
	UpdateFinalizedBlock(Block) error
}
