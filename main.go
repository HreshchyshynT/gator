package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	cfg "github.com/hreshchyshynt/gator/internal/config"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error read config: %v", err)
	}

	state := NewState(config)
	commands := NewCommands()

	commands.Register("login", handleLogin)

	args := os.Args

	if len(args) < 2 {
		log.Fatalln("No arguments provided")
	}

	command := NewCommand(
		args[1],
		args[2:],
	)

	err = commands.Run(&state, command)
	if err != nil {
		log.Fatalf("Can not run %v: %v\n", command.name, err)
	}

	config, err = cfg.Read()
	if err != nil {
		log.Fatalf("error read new config: %v", err)
	}
	fmt.Printf("New config read: %v", config)
}

func handleLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command requires username argument")
	}

	username := cmd.args[0]

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User has been set: %v\n", username)

	return nil
}
