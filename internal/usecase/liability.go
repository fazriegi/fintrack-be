package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/repository"
	"github.com/fazriegi/fintrack-be/pkg"
	"github.com/fazriegi/fintrack-be/pkg/constant"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type liabilityUsecase struct {
	db   *sqlx.DB
	log  *log.Logger
	repo repository.LiabilityRepository
}

type LiabilityUsecase interface {
	ListCategory(ctx context.Context) (resp pkg.Response)
}

func NewLiabilityUsecase(db *sqlx.DB, log *log.Logger, repo repository.LiabilityRepository) LiabilityUsecase {
	return &liabilityUsecase{db, log, repo}
}

func (u *liabilityUsecase) ListCategory(ctx context.Context) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)

	categories, err := u.repo.ListCategory(ctx, userId, u.db)
	if err != nil {
		u.log.Printf("[ERROR] repo.ListCategory: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
	}

	return pkg.NewResponse(http.StatusOK, "Success", categories, nil)
}
