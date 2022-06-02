package repository

import (
	"context"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

func (b *Repository) SaveValidator(validator harvester.Validator) error {
	return b.saveValidator(validator)
}

func (b *Repository) saveValidator(validator harvester.Validator) error {
	_, err := b.ent.Validator.Create().SetAccountID(validator.AccountID).SetName(validator.Name).SetCommission(validator.Commission).SetStatus(validator.Status).SetBalance(validator.Balance).SetReserved(validator.Reserved).SetLocked(validator.Locked).SetOwnStake(validator.OwnStake).SetTotalStake(validator.TotalStake).SetIdentity(validator.Identity).SetNominators(validator.Nominators).SetParent(validator.Parent).SetChildren(validator.Children).SetHash(validator.Hash).SetChain(validator.Chain).Save(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (b *Repository) GetValidatorMap() (map[string]string, error) {
	validators, err := b.ent.Validator.Query().Select("account_id", "hash").All(context.Background())
	if err != nil {
		return nil, err
	}

	validatorMap := make(map[string]string)
	for _, validator := range validators {
		validatorMap[validator.AccountID] = validator.Hash
	}

	return validatorMap, nil
}
