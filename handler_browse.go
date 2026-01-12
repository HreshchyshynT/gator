package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hreshchyshynt/gator/internal/database"
)

func handleBrowse(s *State, command Command, user database.User) error {
	limit := 2
	if len(command.args) > 0 {
		var err error
		limit, err = strconv.Atoi(command.args[0])
		if err != nil {
			return fmt.Errorf("error parsing argument: %v", err)
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Lim:    int32(limit),
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error get posts from db: %v", err)
	}

	for i, p := range posts {
		fmt.Printf("%v. %v\n\n", i+1, p)
	}

	return nil
}
