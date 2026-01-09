package database

import (
	"fmt"
	"time"
)

func (u User) String() string {
	// ID        uuid.UUID
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// Name      string
	return fmt.Sprintf(
		"User(\n\tID: %v\n\tCreatedAt: %v\n\tUpdatedAt: %v\n\tName: %v\n)",
		u.ID,
		u.CreatedAt.Format(time.RFC3339),
		u.UpdatedAt.Format(time.RFC3339),
		u.Name,
	)
}

func (f Feed) String() string {
	return fmt.Sprintf(
		"Feed(\n\tID: %v\n\tCreatedAt: %v\n\tUpdatedAt: %v\n\tName: %v\n\tUrl: %v\n\tUserId: %v\n)",
		f.ID,
		f.CreatedAt.Format(time.RFC3339),
		f.UpdatedAt.Format(time.RFC3339),
		f.Name,
		f.Url,
		f.UserID,
	)
}

func (p Post) String() string {
	return fmt.Sprintf(
		"Post(\nID: %v\n CreatedAt: %v\n UpdatedAt: %v\n Title: %v\n Url: %v\n Description: %v\n PublishedAt: %v\n FeedID: %v\n)",
		p.ID,
		p.CreatedAt.Format(time.RFC3339),
		p.UpdatedAt.Format(time.RFC3339),
		p.Title,
		p.Url,
		p.Description,
		p.PublishedAt.Format(time.RFC3339),
		p.FeedID,
	)
}
