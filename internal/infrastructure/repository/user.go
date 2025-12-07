package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByUsername(username string, db *sqlx.DB) (entity.User, error)
	Insert(data *entity.User, db *sqlx.Tx) (result uint, err error)
}

type userRepo struct {
}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) GetByUsername(username string, db *sqlx.DB) (result entity.User, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.From("users").Select(goqu.I("username"), goqu.I("password"), goqu.I("email"), goqu.I("id"), goqu.I("name")).Where(goqu.I("username").Eq(username))

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

func (r *userRepo) Insert(data *entity.User, tx *sqlx.Tx) (result uint, err error) {
	dialect := pkg.GetDialect()

	dataset := dialect.Insert("users").Rows(*data)
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
