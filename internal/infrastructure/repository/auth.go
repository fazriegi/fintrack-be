package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	InsertRefreshToken(data entity.RefreshToken, db *sqlx.Tx) (result uint, err error)
	GetRefreshToken(token string, userId uint, db *sqlx.DB) (result entity.RefreshToken, err error)
	DeleteRefreshTokenById(id uint, tx *sqlx.Tx) error
	DeleteRefreshTokenByToken(token string, tx *sqlx.Tx) error
	InitUserAssetCategories(userId uint, tx *sqlx.Tx) error
}

type authRepo struct {
}

func NewAuthRepository() AuthRepository {
	return &authRepo{}
}

func (r *authRepo) InsertRefreshToken(data entity.RefreshToken, tx *sqlx.Tx) (result uint, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.Insert("refresh_tokens").Rows(data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	res, err := tx.Exec(sql, val...)
	if err != nil {
		return result, fmt.Errorf("failed to execute insert: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return result, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return uint(id), nil
}

func (r *authRepo) GetRefreshToken(token string, userId uint, db *sqlx.DB) (result entity.RefreshToken, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.From("refresh_tokens").
		Select(
			goqu.I("id"),
			goqu.I("user_id"),
			goqu.I("token"),
			goqu.I("expired_at"),
		).
		Where(
			goqu.I("token").Eq(token),
			goqu.I("user_id").Eq(userId),
		)

	query, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = db.Get(&result, query, val...)
	if err != nil {
		return result, err
	}

	return
}

func (r *authRepo) DeleteRefreshTokenById(id uint, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.Delete("refresh_tokens").Where(goqu.I("id").Eq(id))
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

func (r *authRepo) DeleteRefreshTokenByToken(token string, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.Delete("refresh_tokens").Where(goqu.I("token").Eq(token))
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

func (r *authRepo) InitUserAssetCategories(userId uint, tx *sqlx.Tx) error {
	dialect := pkg.GetDialect()

	dataset := dialect.Insert("user_asset_categories").
		Cols("name", "user_id").
		FromQuery(
			dialect.
				From("asset_categories").
				Select(goqu.I("name"), goqu.L("?", userId)),
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
