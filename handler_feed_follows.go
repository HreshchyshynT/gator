package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
)

func handleFeedFollow(s *State, command Command, user database.User) error {
	if len(command.args) == 0 {
		return fmt.Errorf("follow command required url argument")
	}

	url := command.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Can't find feed for url: %v", url)
	}

	now := time.Now()

	params := database.CreateFeedFollowParams{
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		ID:        uuid.New(),
	}
	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return fmt.Errorf("Error following the feed: %v", err)
	}

	fmt.Printf("User %v now follows %v\n", feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func handleFollowing(s *State, command Command, user database.User) error {
	results, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("Error getting following feeds: %v", err)
	}

	fmt.Printf("User %v is following:\n", user.Name)

	for i, r := range results {
		fmt.Printf("\t%v. %v\n", i+1, r.FeedName)
	}

	return nil
}

func handleUnfollowing(s *State, command Command, user database.User) error {
	if len(command.args) == 0 {
		return fmt.Errorf("unfollow requires url")
	}

	url := command.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)

	if err != nil {
		return fmt.Errorf("Error unfollow feed: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Can not delete feed: %v", err)
	}

	return nil
}
