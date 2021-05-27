package order

import (
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Order represents order
type Order struct {
	ID          uuid.UUID
	Description string
	Book        book.Book
}
