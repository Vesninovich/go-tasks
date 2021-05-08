package book

import (
	"time"

	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Author is author data
type Author struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
