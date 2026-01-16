package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
)

const (
	// is enum better here?
	sortNewest = "newest"
	sortOldest = "oldest"
	sortTitle  = "title"
	sortFeed   = "feed"

	limitArg  = "limit"
	sortArg   = "sort"
	feedArg   = "feed"
	sinceArg  = "since"
	offsetArg = "offset"
	pageArg   = "page"
)

type browseOptions struct {
	limit  int32
	sort   string
	feed   sql.NullString
	since  sql.NullTime
	offset int32
	page   int32
}

func defaultOptions() browseOptions {
	return browseOptions{
		limit:  int32(2),
		feed:   sql.NullString{Valid: false},
		since:  sql.NullTime{Valid: false},
		offset: int32(0),
		page:   int32(0),
	}
}

func handleBrowse(s *State, command Command, user database.User) error {
	options, err := parseArguments(command.Args)

	if err != nil {
		return err
	}

	browse := getBrowseWithSorting(s, options.sort)

	for {
		posts, err := browse(context.Background(), options, user.ID)

		if err != nil {
			return fmt.Errorf("Error get posts from db: %v", err)
		}

		if options.page > 0 {
			fmt.Printf("Displaying results for page: %v\n", options.page)
		}

		for i, p := range posts {
			publishedAt := p.PublishedAt.Format(time.DateTime)
			fmt.Printf(
				"%v. %v\nDescription: %v\nPublishedAt: %v\nUrl: %v\n",
				i+1+int(options.offset),
				p.Title,
				p.Description,
				publishedAt,
				p.Url,
			)
		}
		if len(posts) < int(options.limit) {
			break
		}

		fmt.Println("Press Enter for more...")
		var char rune
		// TODO: fix when typing exit, e is read and xit is passed to terminal
		_, err = fmt.Scanf("%c", &char)
		if err != nil {
			return err
		}
		if char != '\n' {
			break
		}
		options.page += 1
		options.offset += options.limit
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
				Off:      options.offset,
			})
		case sortFeed:
			return s.db.GetPostsForUserFeed(ctx, database.GetPostsForUserFeedParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
				Off:      options.offset,
			})
		case sortTitle:
			return s.db.GetPostsForUserTitle(ctx, database.GetPostsForUserTitleParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
				Off:      options.offset,
			})
		case sortNewest:
			fallthrough
		default:
			return s.db.GetPostsForUserNewest(ctx, database.GetPostsForUserNewestParams{
				UserID:   id,
				Lim:      options.limit,
				Since:    options.since,
				FeedName: options.feed,
				Off:      options.offset,
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
		case name == offsetArg:
			offset, err := strconv.Atoi(arg.Value)
			if err != nil {
				return options, fmt.Errorf("error parsing offset: %v", err)
			}
			options.offset = int32(offset)
		case name == pageArg:
			page, err := strconv.Atoi(arg.Value)
			if err != nil {
				return options, fmt.Errorf("error parsing page: %v", err)
			}
			options.page = int32(page)
		}
	}

	if options.offset > 0 && options.page > 0 {
		fmt.Println("Ignoring offset, page provided.")
	}

	if options.page > 0 {
		options.offset = options.limit * (options.page - 1)
	}

	return options, nil
}
