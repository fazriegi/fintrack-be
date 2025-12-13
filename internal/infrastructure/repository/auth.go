package repository

import (
	"fmt"

	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	InsertRefreshToken(data map[string]any, db *sqlx.Tx) (result uint, err error)
}

type authRepo struct {
}

func NewAuthRepository() AuthRepository {
	return &authRepo{}
}

func (r *authRepo) InsertRefreshToken(data map[string]any, tx *sqlx.Tx) (result uint, err error) {
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
