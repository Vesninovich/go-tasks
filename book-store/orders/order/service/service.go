package service

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	catalogservice "github.com/Vesninovich/go-tasks/book-store/orders/catalog/service"
	"github.com/Vesninovich/go-tasks/book-store/orders/order"
)

// Service handles orders manipulation
type Service struct {
	repo    order.Repository
	catalog *catalogservice.Service
}

// New creates new instance of Service
func New(r order.Repository, c *catalogservice.Service) *Service {
	return &Service{r, c}
}

// GetOrder reads stored order by id
func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (order.Order, error) {
	var empty order.Order
	var book book.Book
	dto, err := s.repo.Get(ctx, id)
	if err != nil {
		return empty, err
	}
	book, err = s.catalog.GetBook(ctx, dto.BookID)
	if err != nil {
		return empty, err
	}
	return order.Order{
		ID:          dto.ID,
		Description: dto.Description,
		Book:        book,
	}, err
}

// CreateOrder validates data, creates order if data is valid and saves it, returns error otherwise.
func (s *Service) CreateOrder(ctx context.Context, data order.CreateDTO) (order.Order, error) {
	var empty order.Order
	if data.Description == "" {
		return empty, &commonerrors.InvalidInput{Reason: "Description is required"}
	}
	if data.BookID.IsZero() {
		return empty, &commonerrors.InvalidInput{Reason: "Book ID is required"}
	}
	b, err := s.catalog.GetBook(ctx, data.BookID)
	if err != nil {
		return empty, err
	}
	res, err := s.repo.Create(ctx, data)
	if err != nil {
		return empty, err
	}
	return order.Order{
		ID:          res.ID,
		Description: res.Description,
		Book:        b,
	}, err
}
