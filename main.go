package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	cfg "github.com/hreshchyshynt/gator/internal/config"
	"github.com/hreshchyshynt/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error read config: %v", err)
	}

	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		log.Fatalln("Can't open database connection")
	}

	state := NewState(config, database.New(db))
	commands := NewCommands()

	commands.Register("login", handleLogin)
	commands.Register("register", handleRegister)

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
}

func handleLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command requires username argument")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("User has been set: %v\n", username)

	return nil
}

func handleRegister(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("Register command requires username argument")
	}

	username := cmd.args[0]

	now := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), params)

	if err != nil {
		return err
	}

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User has been created: %v\n", user)

	return nil
}
