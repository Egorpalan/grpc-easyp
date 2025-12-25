package notes

import "time"

type Note struct {
	ID          *string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
