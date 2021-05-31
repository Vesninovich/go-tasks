package grpc

import (
	"context"

	"github.com/Vesninovich/go-tasks/book-store/common/orders"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/Vesninovich/go-tasks/book-store/orders/order"
	orderservice "github.com/Vesninovich/go-tasks/book-store/orders/order/service"
)

// Server of orders
type Server struct {
	orders.UnimplementedOrdersServer

	service *orderservice.Service
}

// New creates Server
func New(s *orderservice.Service) *Server {
	return &Server{
		service: s,
	}
}

// GetOrder godoc
func (s *Server) GetOrder(ctx context.Context, req *orders.ID) (*orders.Order, error) {
	id, err := uuid.FromBytes(req.Id)
	if err != nil {
		return nil, err
	}
	o, err := s.service.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	return orderToResponse(o), nil
}

// CreateOrder godoc
func (s *Server) CreateOrder(ctx context.Context, dto *orders.CreateDTO) (*orders.Order, error) {
	bID, err := uuid.FromBytes(dto.Book)
	if err != nil {
		return nil, err
	}
	o, err := s.service.CreateOrder(ctx, order.CreateDTO{
		Description: dto.Description,
		BookID:      bID,
	})
	if err != nil {
		return nil, err
	}
	return orderToResponse(o), nil
}

// UpdateDescription godoc
func (s *Server) UpdateDescription(ctx context.Context, dto *orders.DescriptionUpdate) (*orders.Order, error) {
	id, err := uuid.FromBytes(dto.Id)
	if err != nil {
		return nil, err
	}
	o, err := s.service.UpdateDescription(ctx, order.Order{
		ID:          id,
		Description: dto.Description,
	})
	if err != nil {
		return nil, err
	}
	return orderToResponse(o), nil
}

// RemoveOrder godoc
func (s *Server) RemoveOrder(ctx context.Context, id *orders.ID) (*orders.Order, error) {
	oid, err := uuid.FromBytes(id.Id)
	if err != nil {
		return nil, err
	}
	o, err := s.service.RemoveOrder(ctx, oid)
	if err != nil {
		return nil, err
	}
	return orderToResponse(o), nil
}

func orderToResponse(o order.Order) *orders.Order {
	return &orders.Order{
		Id:          o.ID[:],
		Description: o.Description,
		Book:        o.Book.ID[:],
	}
}
