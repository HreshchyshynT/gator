package main

import (
	"context"
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

	argLimit = "limit"
	argSort  = "sort"
)

func handleBrowse(s *State, command Command, user database.User) error {
	limit := 2

	argsMap, err := parseArguments(command.Args)

	if err != nil {
		return err
	}

	if _, ok := argsMap[argLimit]; ok {
		limit, _ = strconv.Atoi(argsMap[argLimit])
	}

	browse := getBrowseWithSorting(s, argsMap[argSort])

	posts, err := browse(context.Background(), int32(limit), user.ID)

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
) func(context.Context, int32, uuid.UUID) ([]database.Post, error) {
	return func(ctx context.Context, limit int32, id uuid.UUID) ([]database.Post, error) {
		switch strings.ToLower(sort) {
		case sortOldest:
			return s.db.GetPostsForUserOldest(ctx, database.GetPostsForUserOldestParams{
				UserID: id,
				Lim:    limit,
			})
		case sortFeed:
			return s.db.GetPostsForUserFeed(ctx, database.GetPostsForUserFeedParams{
				UserID: id,
				Lim:    limit,
			})
		case sortTitle:
			return s.db.GetPostsForUserTitle(ctx, database.GetPostsForUserTitleParams{
				UserID: id,
				Lim:    limit,
			})
		case sortNewest:
			fallthrough
		default:
			return s.db.GetPostsForUserNewest(ctx, database.GetPostsForUserNewestParams{
				UserID: id,
				Lim:    limit,
			})
		}
	}
}

func parseArguments(args []Argument) (map[string]string, error) {
	result := make(map[string]string)
	if len(args) == 0 {
		return result, nil
	}

	// handle case when limit passed as first unnamed parameter
	if len(args[0].Name) == 0 {
		result[argLimit] = args[0].Value
	}

	// put every named argument to map result
	// but only required by this handler will be used
	// probably it will be common logic for most commands
	for _, arg := range args {
		if len(arg.Name) > 0 {
			result[arg.Name] = arg.Value
		}
	}

	// validate limit
	if limit, ok := result[argLimit]; ok {
		_, err := strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("error parsing limit: %v", err)
		}
	}

	return result, nil
}
