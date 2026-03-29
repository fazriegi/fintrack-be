package domain

import "github.com/google/uuid"

type Category struct {
	Id       uuid.UUID `db:"id" json:"id"`
	UserID   uuid.UUID `db:"user_id" json:"-"`
	Name     string    `db:"name" json:"name"`
	BaseType string    `db:"base_type" json:"base_type"`
}
