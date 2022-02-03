package order

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/stored"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// CreateDTO is DTO for creating order
type CreateDTO struct {
	Description string
	BookID      uuid.UUID
}

// DTO is DTO of order
type DTO struct {
	ID uuid.UUID
	CreateDTO
}

// StoredOrderDTO is order that is stored
type StoredOrderDTO struct {
	DTO
	stored.Stored
}

// Repository of orders
type Repository interface {
	GetAll(ctx context.Context) ([]DTO, error)
	Get(ctx context.Context, id uuid.UUID) (DTO, error)
	Create(ctx context.Context, dto CreateDTO) (DTO, error)
	Update(ctx context.Context, dto DTO) (DTO, error)
	Delete(ctx context.Context, id uuid.UUID) (DTO, error)
}

// ToOrderDTO converts stored version to general DTO
func (s StoredOrderDTO) ToOrderDTO() DTO {
	return DTO{
		ID: s.ID,
		CreateDTO: CreateDTO{
			Description: s.Description,
			BookID:      s.BookID,
		},
	}
}
