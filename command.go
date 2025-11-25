package main

import (
	"fmt"
)

type CommandCallback func(*State, Command) error

type Command struct {
	name string
	args []string
}

func NewCommand(name string, args []string) Command {
	return Command{
		name: name,
		args: args,
	}
}

type Commands struct {
	registry map[string]CommandCallback
}

func (c *Commands) Run(s *State, cmd Command) error {
	callback, ok := c.registry[cmd.name]
	if !ok {
		return fmt.Errorf("Command %v not found", cmd.name)
	}
	return callback(s, cmd)
}

func (c *Commands) Register(name string, callback CommandCallback) {
	c.registry[name] = callback
}

func NewCommands() Commands {
	return Commands{
		registry: make(map[string]CommandCallback),
	}
}
