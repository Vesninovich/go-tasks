package book

import (
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Author is author data
type Author struct {
	ID   uuid.UUID
	Name string
}
