package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
)

func handleLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("login command requires username argument")
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
		return errors.New("register command requires username argument")
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

func handleGetAllUsers(s *State, cmd Command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("can't get all users: %v", err)
	}

	var buffer strings.Builder

	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			buffer.WriteString("* " + u.Name + " (current)\n")
		} else {
			buffer.WriteString("* " + u.Name + "\n")
		}
	}

	fmt.Print(buffer.String())

	return nil
}
