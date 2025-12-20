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
