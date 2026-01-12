package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
	"github.com/hreshchyshynt/gator/internal/rss"
)

func handleListFeeds(s *State, command Command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Can't retrieve feeds: %v", err)
	}

	var builder strings.Builder

	for i, f := range feeds {
		builder.WriteString(strconv.Itoa(i+1) + ": ")
		builder.WriteString(f.Name + ", " + f.Url + ", added by " + f.Username + "\n")
	}

	fmt.Print(builder.String())

	return nil
}

func handleAddFeed(s *State, command Command, user database.User) error {
	args := command.args
	if len(args) < 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		return fmt.Errorf("addfeed command requires not empty feed name and url arguments")
	}

	feedName := args[0]
	feedUrl := args[1]

	_, err := rss.FetchFeed(context.Background(), feedUrl)

	if err != nil {
		return err
	}

	now := time.Now()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("Can't save feed to db: %v", err)
	}

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	})

	fmt.Printf("Feed receive: %v\n", feed)

	return nil
}
