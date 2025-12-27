package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type AssetRepository interface {
	ListAssetCategory(userId uint, db *sqlx.DB) ([]entity.AssetCategory, error)
	ListAsset(param entity.ListAssetRequest, db *sqlx.DB) ([]entity.AssetResponse, uint, error)
	Insert(data entity.Asset, tx *sqlx.Tx) error
	GetById(id, userId uint, forUpdate bool, db *sqlx.DB, tx *sqlx.Tx) (result entity.AssetResponse, err error)
	Update(data entity.Asset, id, userId uint, tx *sqlx.Tx) error
	Delete(id, userId uint, tx *sqlx.Tx) error
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

func (r assetRepo) ListAsset(param entity.ListAssetRequest, db *sqlx.DB) (result []entity.AssetResponse, count uint, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.From(goqu.T("assets").As("a")).
		Join(goqu.T("user_asset_categories").As("b"), goqu.On(
			goqu.I("a.category_id").Eq(goqu.I("b.id")),
			goqu.I("a.user_id").Eq(goqu.I("b.user_id")),
		)).
		Select(
			goqu.I("a.id"),
			goqu.I("a.name"),
			goqu.I("b.name").As("category"),
			goqu.I("a.amount"),
			goqu.I("a.purchase_price"),
			goqu.I("a.status"),
			goqu.L("a.amount * a.purchase_price").As("total_purchase_price"),
		).
		Where(
			goqu.I("a.user_id").Eq(param.UserId),
		)

	if param.Name != "" {
		dataset = dataset.Where(goqu.I("a.name").ILike("%" + param.Name + "%"))
	}

	if param.Category != "" {
		dataset = dataset.Where(goqu.I("b.name").ILike("%" + param.Category + "%"))
	}

	if param.Status != "" {
		dataset = dataset.Where(goqu.I("a.status").ILike("%" + param.Status + "%"))
	}

	result = make([]entity.AssetResponse, 0)
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		dataset = pkg.QueryWithPagination(dataset, param.PaginationRequest)

		sql, val, err := dataset.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build SQL query: %w", err)
		}

		row, err := db.Queryx(sql, val...)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		defer row.Close()

		err = pkg.ScanRowsIntoStructs(row, &result)
		if err != nil {
			return fmt.Errorf("failed to scan rows into structs: %w", err)
		}

		query, val, err := dataset.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build SQL query: %w", err)
		}

		err = db.Select(&result, query, val...)
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		datasetCount := dataset.Select(goqu.COUNT("*").As("total"))

		query, val, err := datasetCount.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build SQL query: %w", err)
		}

		err = db.Get(&count, query, val...)
		if err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
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

func (r *assetRepo) GetById(id, userId uint, forUpdate bool, db *sqlx.DB, tx *sqlx.Tx) (result entity.AssetResponse, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.From(goqu.T("assets").As("a")).
		Join(goqu.T("user_asset_categories").As("b"), goqu.On(
			goqu.I("a.category_id").Eq(goqu.I("b.id")),
			goqu.I("a.user_id").Eq(goqu.I("b.user_id")),
		)).
		Select(
			goqu.I("a.id"),
			goqu.I("a.name"),
			goqu.I("a.category_id"),
			goqu.I("b.name").As("category"),
			goqu.I("a.amount"),
			goqu.I("a.purchase_price"),
			goqu.I("a.status"),
			goqu.L("a.amount * a.purchase_price").As("total_purchase_price"),
		).
		Where(
			goqu.I("a.user_id").Eq(userId),
			goqu.I("a.id").Eq(id),
		)

	if tx != nil && forUpdate {
		dataset = dataset.ForUpdate(exp.Wait)
	}

	query, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	if tx != nil && forUpdate {
		err = tx.Get(&result, query, val...)
		if err != nil {
			return result, err
		}
	} else {
		err = db.Get(&result, query, val...)
		if err != nil {
			return result, err
		}
	}

	return
}

func (r *assetRepo) Update(data entity.Asset, id, userId uint, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.
		Update("assets").
		Set(data).
		Where(
			goqu.I("user_id").Eq(userId),
			goqu.I("id").Eq(id),
		)

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

func (r *assetRepo) Delete(id, userId uint, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.
		Delete("assets").
		Where(
			goqu.I("user_id").Eq(userId),
			goqu.I("id").Eq(id),
		)

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
