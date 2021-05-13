package book

import (
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Book is book data
type Book struct {
	ID         uuid.UUID
	Name       string
	Author     Author
	Categories []Category
}
