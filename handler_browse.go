package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
)

const (
	// is enum better here?
	sortNewest = "newest"
	sortOldest = "oldest"
	sortTitle  = "title"
	sortFeed   = "feed"

	limitArg = "limit"
	sortArg  = "sort"
	feedArg  = "feed"
	sinceArg = "since"
)

type browseOptions struct {
	limit int32
	sort  string
	feed  sql.NullString
	since sql.NullTime
}

func defaultOptions() browseOptions {
	return browseOptions{
		limit: int32(2),
		feed:  sql.NullString{Valid: false},
		since: sql.NullTime{Valid: false},
	}
}

func handleBrowse(s *State, command Command, user database.User) error {
	options, err := parseArguments(command.Args)

	if err != nil {
		return err
	}

	browse := getBrowseWithSorting(s, options.sort)

	posts, err := browse(context.Background(), options, user.ID)

	if err != nil {
		return fmt.Errorf("Error get posts from db: %v", err)
	}

	for i, p := range posts {
		fmt.Printf("%v. %v\n\n", i+1, p)
	}

	return nil
}

func getBrowseWithSorting(
	s *State,
	sort string,
) func(context.Context, browseOptions, uuid.UUID) ([]database.Post, error) {
	return func(ctx context.Context, options browseOptions, id uuid.UUID) ([]database.Post, error) {
		switch strings.ToLower(sort) {
		case sortOldest:
			return s.db.GetPostsForUserOldest(ctx, database.GetPostsForUserOldestParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
			})
		case sortFeed:
			return s.db.GetPostsForUserFeed(ctx, database.GetPostsForUserFeedParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
			})
		case sortTitle:
			return s.db.GetPostsForUserTitle(ctx, database.GetPostsForUserTitleParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
			})
		case sortNewest:
			fallthrough
		default:
			return s.db.GetPostsForUserNewest(ctx, database.GetPostsForUserNewestParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
			})
		}
	}
}

func parseArguments(args []Argument) (browseOptions, error) {
	options := defaultOptions()

	if len(args) == 0 {
		return options, nil
	}

	for i, arg := range args {
		name := arg.Name
		switch {
		case name == limitArg || (i == 0 && len(arg.Name) == 0):
			limit, err := strconv.Atoi(arg.Value)
			if err != nil {
				return options, fmt.Errorf("error parsing limit: %v", err)
			}
			options.limit = int32(limit)

		case name == sortArg:
			options.sort = arg.Value
		case name == feedArg:
			options.feed = sql.NullString{
				String: arg.Value,
				Valid:  len(arg.Value) > 0,
			}
		case name == sinceArg:
			time, err := parsePostDate(arg.Value)
			if err != nil {
				fmt.Printf("%v, ignoring \"since\" parameter\n", err)
			}
			options.since = sql.NullTime{
				Time:  time,
				Valid: err == nil,
			}
		}
	}

	return options, nil
}
