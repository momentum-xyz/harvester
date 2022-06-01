package harvester

type Repository interface {
	GetBlock(number uint64) (Block, error)
	GetValidatorMap() (map[string]string, error)
	SaveBlock(Block) error
	SaveValidator(validator Validator) error
	UpdateFinalizedBlock(Block) error
}
