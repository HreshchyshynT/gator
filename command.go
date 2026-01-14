package main

import (
	"fmt"
	"strings"
)

type CommandCallback func(*State, Command) error

type Argument struct {
	Name  string
	Value string
}

func (a Argument) HasName() bool {
	return len(a.Name) > 0
}

type Command struct {
	Name string
	Args []Argument
}

func (c Command) HasArgument(name string) bool {
	for _, arg := range c.Args {
		if arg.Name == name {
			return true
		}
	}
	return false
}

func (c Command) String() string {
	var builder strings.Builder

	builder.WriteString(c.Name + " ")
	for _, arg := range c.Args {
		if arg.HasName() {
			builder.WriteString("--" + arg.Name + "=")
		}
		builder.WriteString(arg.Value + " ")
	}

	return builder.String()
}

func NewCommand(name string, args []string) Command {
	var cmdArgs []Argument
	for _, arg := range args {
		var value string
		var argName string
		if len(arg) > 2 && arg[:2] == "--" {
			if splitted := strings.SplitN(arg, "=", 2); len(splitted) > 1 {
				argName = splitted[0][2:]
				value = splitted[1]
			} else {
				value = arg[2:]
			}
		} else {
			value = arg
		}

		cmdArgs = append(cmdArgs, Argument{
			Name:  argName,
			Value: value,
		})
	}
	return Command{
		Name: name,
		Args: cmdArgs,
	}
}

type Commands struct {
	registry map[string]CommandCallback
}

func (c *Commands) Run(s *State, cmd Command) error {
	callback, ok := c.registry[cmd.Name]
	if !ok {
		return fmt.Errorf("Command %v not found", cmd.Name)
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
