package cron

import (
	"log"

	"github.com/fazriegi/fintrack-be/internal/repository"
	"github.com/jmoiron/sqlx"
)

func Start(db *sqlx.DB, logger *log.Logger) {
	// Refresh Token
	userRepo := repository.NewUserRepository()
	go func() {
		RefreshTokenCleanup(db, userRepo, logger)
	}()

	// Networth
	networthRepo := repository.NewNetworthRepository()
	go func() {
		NetworthCalculate(db, networthRepo, logger)
	}()
}
