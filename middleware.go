package main

import (
	"context"
	"fmt"

	"github.com/hreshchyshynt/gator/internal/database"
)

func middlewareLoggedIn(handler func(*State, Command, database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		userName := s.config.CurrentUserName
		user, err := s.db.GetUser(context.Background(), userName)
		if err != nil {
			return fmt.Errorf("Can't get current user: %v", err)
		}
		return handler(s, cmd, user)
	}
}
