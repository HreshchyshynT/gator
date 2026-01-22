package rss

import (
	"encoding/xml"
	"fmt"
)

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

func (i *RSSItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type xmlData struct {
		XMLName xml.Name
		Data    string `xml:",chardata"`
	}

	type raw RSSItem

	type TempItem struct {
		raw
		Links []xmlData `xml:"link"`
	}

	var temp TempItem

	err := dec.DecodeElement(&temp, &start)
	if err != nil {
		return err
	}

	i.Title = temp.Title
	i.Description = temp.Description
	i.PubDate = temp.PubDate
	for _, l := range temp.Links {
		if l.XMLName.Space == "" {
			i.Link = l.Data
			break
		}
	}
	return nil
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
