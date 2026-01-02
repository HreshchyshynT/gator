package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
	"github.com/hreshchyshynt/gator/internal/rss"
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

func handleReset(s *State, _ Command) error {
	err := s.db.ClearUsers(context.Background())
	if err == nil {
		fmt.Println("Table \"users\" cleared.")
	}
	return err
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

func handleAggregate(s *State, command Command) error {
	if len(command.args) == 0 || len(command.args[0]) == 0 {
		return fmt.Errorf("agg command requires not empty feed url argument")
	}
	feed, err := rss.FetchFeed(context.Background(), command.args[0])
	if err != nil {
		return err
	}
	fmt.Printf(
		"feed received: %v, link (%v), description (%v), items count: %v\n",
		feed.Channel.Title,
		feed.Channel.Link,
		feed.Channel.Description,
		len(feed.Channel.Items),
	)
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
