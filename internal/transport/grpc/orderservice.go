package grpc

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"lyceum_service/internal/models"
	client "lyceum_service/pkg/api/order"
)

type Service interface {
	CreatePosition(ctx context.Context, position models.Position) (*models.Position, error)
}

type OrderService struct {
	client.UnimplementedOrderServiceServer
	service Service
}

func NewOrderService(srv Service) *OrderService {
	return &OrderService{service: srv}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *client.CreateOrderRequest) (*client.CreateOrderResponse, error) {
	return &client.CreateOrderResponse{}, nil
}

func (s *OrderService) CreatePosition(ctx context.Context, req *client.CreatePositionRequest) (*client.CreatePositionResponse, error) {
	resp, err := s.service.CreatePosition(ctx, models.Position{
		Price: req.GetPrice(),
		Name:  req.GetName(),
	})
	if err != nil {
		return nil, fmt.Errorf("CreatePosition: %w", err)
	}
	r := pointer.Get(resp)
	return &client.CreatePositionResponse{
		Id:    r.ID,
		Name:  r.Name,
		Price: r.Price,
	}, nil
}
