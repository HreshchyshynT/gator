package main

import (
	cfg "github.com/hreshchyshynt/gator/internal/config"
	"github.com/hreshchyshynt/gator/internal/database"
)

type State struct {
	db     *database.Queries
	config *cfg.Config
}

func NewState(
	config *cfg.Config,
	db *database.Queries,
) State {
	return State{
		config: config,
		db:     db,
	}
}
