package book

import (
	"time"

	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Book is book data
type Book struct {
	ID         uuid.UUID
	Name       string
	Author     Author
	Categories []Category
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
