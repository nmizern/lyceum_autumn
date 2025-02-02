package service

import (
	"context"
	"lyceum_service/internal/models"
)

type OrderRepo interface {
	CreatePosition(ctx context.Context, position models.Position) (*models.Position, error)
}

type OrderService struct {
	Repo OrderRepo
}

func NewOrderService(repo OrderRepo) *OrderService {
	return &OrderService{repo}
}

func (s *OrderService) CreatePosition(ctx context.Context, position models.Position) (*models.Position, error) {
	return s.Repo.CreatePosition(ctx, position)
}
