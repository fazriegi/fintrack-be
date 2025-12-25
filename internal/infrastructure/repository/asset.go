package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
)

type AssetRepository interface {
	ListAssetCategory(userId uint, db *sqlx.DB) ([]entity.AssetCategory, error)
	Insert(data entity.Asset, tx *sqlx.Tx) error
}

type assetRepo struct {
}

func NewAssetRepository() AssetRepository {
	return &assetRepo{}
}

func (r assetRepo) ListAssetCategory(userId uint, db *sqlx.DB) (result []entity.AssetCategory, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.From("user_asset_categories").
		Select(
			goqu.I("id"),
			goqu.I("name"),
		).
		Where(
			goqu.I("user_id").Eq(userId),
		)

	query, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = db.Select(&result, query, val...)
	if err != nil {
		return result, err
	}

	return
}

func (r *assetRepo) Insert(data entity.Asset, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.Insert("assets").Rows(data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = tx.Exec(sql, val...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}
