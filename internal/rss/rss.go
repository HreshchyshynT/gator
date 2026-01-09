package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedUrl string) (RSSFeed, error) {
	var feed RSSFeed
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, http.NoBody)

	if err != nil {
		return feed, fmt.Errorf("Can't create request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return feed, fmt.Errorf("Can't perform request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return feed, fmt.Errorf("Error fetching feed: response code %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return feed, fmt.Errorf("Can't read response body: %v", err)
	}

	if err := xml.Unmarshal(data, &feed); err != nil {
		return feed, fmt.Errorf("Can't unmarshal response body: %v", err)
	}

	unescape(&feed)

	return feed, nil
}

func unescape(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}
