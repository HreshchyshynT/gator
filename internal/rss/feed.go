package rss

import "fmt"

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (ri RSSItem) String() string {
	return fmt.Sprintf(
		"RSSItem(\n\ttitle: %v\n\tlink: %v\n\tdescription: %v\n\tpubDate: %v\n)",
		ri.Title,
		ri.Link,
		ri.Description,
		ri.PubDate,
	)
}
