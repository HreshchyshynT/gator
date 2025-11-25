package main

import cfg "github.com/hreshchyshynt/gator/internal/config"

type State struct {
	config *cfg.Config
}

func NewState(config *cfg.Config) State {
	return State{
		config: config,
	}
}
