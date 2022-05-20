package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

// Block holds the schema definition for the Block entity.
type Block struct {
	ent.Schema
}

// Fields of the Block.
func (Block) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("number").Positive(),
		field.String("author_id"),
		field.Bool("finalized"),
		field.Int("extrinsics_count"),
		field.JSON("extrinsics", []harvester.Extrinsic{}).Optional(),
		field.String("chain"),
	}
}

// Edges of the Block.
func (Block) Edges() []ent.Edge {
	return nil
}
