package repository

import (
	"context"

	"github.com/fazriegi/fintrack-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type liabilityRepository struct{}

type LiabilityRepository interface {
	ListCategory(ctx context.Context, userId uuid.UUID, db *sqlx.DB) (*[]domain.Category, error)
}

func NewLiabilityRepository() LiabilityRepository {
	return &liabilityRepository{}
}

func (r *liabilityRepository) ListCategory(ctx context.Context, userId uuid.UUID, db *sqlx.DB) (*[]domain.Category, error) {
	var categories = make([]domain.Category, 0)
	query := `SELECT id, name, base_type FROM liability_categories WHERE user_id = $1 ORDER BY name ASC`
	err := db.SelectContext(ctx, &categories, query, userId)

	return &categories, err
}
