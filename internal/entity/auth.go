package entity

import "time"

type (
	RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RefreshToken struct {
		ID        uint      `db:"id"`
		UserId    uint      `db:"user_id" `
		Token     string    `db:"token" `
		ExpiredAt time.Time `db:"expired_at" `
	}
)
