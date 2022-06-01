package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

// Validator holds the schema definition for the Validator entity.
type Validator struct {
	ent.Schema
}

// Fields of the Validator.
func (Validator) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id"),
		field.String("name"),
		field.Float("commission"),
		field.String("status"),
		field.String("balance"),
		field.String("reserved"),
		field.JSON("locked", []harvester.ValidatorBalancesLocked{}),
		field.String("own_stake"),
		field.String("total_stake"),
		field.JSON("identity", harvester.ValidatorInfo{}),
		field.JSON("nominators", []harvester.Nominator{}),
		field.JSON("parent", harvester.Parent{}),
		field.JSON("children", []string{}),
		field.String("hash"),
		field.String("chain"),
	}
}

// Edges of the Validator.
func (Validator) Edges() []ent.Edge {
	return nil
}
