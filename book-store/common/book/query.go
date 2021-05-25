package book

import "github.com/Vesninovich/go-tasks/book-store/common/uuid"

// Query represents query for books
type Query struct {
	ID         uuid.UUID
	Author     uuid.UUID
	Categories []uuid.UUID
}
