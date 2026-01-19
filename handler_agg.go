package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/gator/internal/database"
	"github.com/hreshchyshynt/gator/internal/rss"
	"github.com/lib/pq"
)

func handleAggregate(s *State, command Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("agg command requires time_between_reqs param")
	}

	timeBetweenReqs, err := time.ParseDuration(command.Args[0].Value)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeed(s.db)
		if err != nil {
			return err
		}
	}
}

func scrapeFeed(db *database.Queries) error {
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Can not get next feed for fetching: %v", err)
	}
	fmt.Printf("Fetching posts for %v...\n", feed.Name)

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Can not fetch feed update: %v", err)
	}

	feed, err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Can not mark feed as fetched: %v", err)
	}

	for _, item := range rssFeed.Channel.Items {
		now := time.Now()
		pubAt, err := parsePostDate(item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing date (%v): %v\n", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})
		var pqErr *pq.Error
		// duplicated keys error has code 23505
		if err != nil && errors.As(err, &pqErr) && pqErr.Code != "23505" {
			fmt.Printf("Error during creating post: %v\n", pqErr)
		}
	}

	fmt.Println("Posts fetched")

	return nil
}
