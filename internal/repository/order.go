package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"lyceum_service/internal/models"
	"lyceum_service/pkg/db/postgres"
)

type OrderRepository struct {
	db *postgres.DB
}

func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (s *OrderRepository) CreatePosition(ctx context.Context, position models.Position) (*models.Position, error) {
	var result models.Position
	err := sq.Insert("position").
		Columns("name", "price").
		Values(position.Name, position.Price).
		Suffix("returning *").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db.Db).
		QueryRow().
		Scan(&result.ID, &result.Name, &result.Price, &result.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.CreatePosition: %w", err)
	}

	return &result, nil
}
