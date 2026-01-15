package main

import (
	"fmt"
	"time"
)

func parsePostDate(dateString string) (t time.Time, err error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC3339,
		"02-01-2006",
		"02/01/2006",
		"01/2006",
		"2006",
	}
	for _, l := range layouts {
		t, err = time.Parse(l, dateString)
		if err == nil {
			return t, err
		}
	}

	return t, fmt.Errorf("Failed to parse date string: %v", dateString)
}
