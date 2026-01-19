package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
	"github.com/hreshchyshynt/gator/internal/utils"
)

const (
	nameArg = "name"
)

type followOptions struct {
	name string
	url  string
}

func (fo followOptions) String() string {
	var builder strings.Builder

	if fo.name != "" {
		builder.WriteString(fmt.Sprintf("name = %v", fo.name))
		if fo.url != "" {
			builder.WriteString(", ")
		}
	}
	if fo.url != "" {
		builder.WriteString(fmt.Sprintf("url = %v", fo.url))
	}

	return builder.String()
}

func newFollowOptions() followOptions {
	return followOptions{}
}

func handleFeedFollow(s *State, command Command, user database.User) error {
	options, err := parseOptions(command.Args)
	if err != nil {
		return err
	}
	dbCall := getFeedCall(options)
	feed, err := dbCall(*s.db, context.Background())

	if err != nil {
		return fmt.Errorf("Can't find feed for %v", options)
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

	if utils.IsDuplicatedKeys(err) {
		return fmt.Errorf("Already following %v", feed.Name)
	}

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
	options, err := parseOptions(command.Args)
	if err != nil {
		return err
	}

	dbCall := getFeedCall(options)
	feed, err := dbCall(*s.db, context.Background())

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

func getFeedCall(
	options followOptions,
) func(database.Queries, context.Context) (database.Feed, error) {
	switch {
	case options.name != "":
		return func(db database.Queries, ctx context.Context) (database.Feed, error) {
			return db.GetFeedByName(ctx, options.name)
		}
	default:
		return func(db database.Queries, ctx context.Context) (database.Feed, error) {
			return db.GetFeedByUrl(ctx, options.url)
		}
	}
}
func parseOptions(args []Argument) (followOptions, error) {
	options := newFollowOptions()
	if len(args) == 0 {
		return options, fmt.Errorf("command required either unnamed feed url or --name=\"feed name\" argument")
	}

	switch args[0].Name {
	case "":
		options.url = args[0].Value
	case nameArg:
		options.name = args[0].Value
	}
	return options, nil
}
