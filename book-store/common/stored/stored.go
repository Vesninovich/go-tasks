package stored

import "time"

// Stored represents storable struct
type Stored struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// IsStored indicates if stored item is already stored
func (s Stored) IsStored() bool {
	return !s.CreatedAt.IsZero()
}

// IsDeleted indicates if stored item is already deleted
func (s Stored) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}
