package main

import (
	"context"
	"fmt"
)

func handleReset(s *State, _ Command) error {
	err := s.db.ClearUsers(context.Background())
	if err == nil {
		fmt.Println("Table \"users\" cleared.")
	}
	return err
}
