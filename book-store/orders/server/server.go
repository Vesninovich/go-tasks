package server

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

// GetOrder gets order
func (s *Server) GetOrder(ctx context.Context, req *orders.OrderID) (*orders.Order, error) {
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

// CreateOrder creates order
func (s *Server) CreateOrder(ctx context.Context, dto *orders.OrderCreateDTO) (*orders.Order, error) {
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

func orderToResponse(o order.Order) *orders.Order {
	return &orders.Order{
		Id:          o.ID[:],
		Description: o.Description,
		Book:        o.Book.ID[:],
	}
}
