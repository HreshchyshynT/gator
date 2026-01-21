package rss

import (
	"encoding/xml"
	"fmt"
	"strings"
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
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "title":
				var v string
				if err := dec.DecodeElement(&v, &t); err != nil {
					return err
				}
				i.Title = strings.TrimSpace(v)

			case "link":
				// Distinguish by namespace:
				if t.Name.Space == "" {
					// RSS2 <link>text</link>
					var v string
					if err := dec.DecodeElement(&v, &t); err != nil {
						return err
					}
					v = strings.TrimSpace(v)
					// Keep the RSS text link (don’t overwrite with empty)
					if v != "" {
						i.Link = v
					}
				} else {
					// some other namespace: just consume it
					var skip any
					if err := dec.DecodeElement(&skip, &t); err != nil {
						return err
					}
				}

			case "description":
				var v string
				if err := dec.DecodeElement(&v, &t); err != nil {
					return err
				}
				i.Description = strings.TrimSpace(v)

			case "pubDate":
				var v string
				if err := dec.DecodeElement(&v, &t); err != nil {
					return err
				}
				i.PubDate = strings.TrimSpace(v)

			default:
				// Consume any element we don't care about (guid, category, media:content, dc:creator, etc.)
				var skip any
				if err := dec.DecodeElement(&skip, &t); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if t.Name.Local == start.Name.Local && t.Name.Space == start.Name.Space {
				return nil
			}
		}
	}
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
