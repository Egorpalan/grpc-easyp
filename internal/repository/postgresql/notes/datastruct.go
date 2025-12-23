package notes

import (
	"time"

	"github.com/Egorpalan/grpc-easyp/internal/model/notes"
)

type Note struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (d *Note) ToServiceModel() *notes.Note {
	return &notes.Note{
		ID:          &d.ID,
		Title:       d.Title,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
