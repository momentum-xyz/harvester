package repository

import (
	"context"
	"database/sql"

	"github.com/OdysseyMomentumExperience/harvester/ent"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
)

type Repository struct {
	DB  *sql.DB
	ent *ent.Client
}

func NewRepository(db *sql.DB, migrate bool) (harvester.Repository, error) {
	ent := NewEnt(db)

	if migrate {
		log.Infof("Migrating database")
		err := ent.Schema.Create(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return &Repository{
		DB:  db,
		ent: ent,
	}, nil
}
