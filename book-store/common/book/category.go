package book

import (
	"time"

	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Category is book category data
type Category struct {
	ID        uuid.UUID
	Name      string
	ParentID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
