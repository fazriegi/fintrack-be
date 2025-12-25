package usecase

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/infrastructure/database"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/repository"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/sirupsen/logrus"
)

type AssetUsecase interface {
	ListAssetCategory(userId uint) (resp pkg.Response)
}

type assetUsecase struct {
	assetRepo repository.AssetRepository
	log       *logrus.Logger
	jwt       *pkg.JWT
}

func NewAssetUsecase(assetRepo repository.AssetRepository, jwt *pkg.JWT) AssetUsecase {
	log := logger.Get()

	return &assetUsecase{
		assetRepo,
		log,
		jwt,
	}
}

func (u *assetUsecase) ListAssetCategory(userId uint) (resp pkg.Response) {
	var (
		err error
		db  = database.Get()
	)

	data, err := u.assetRepo.ListAssetCategory(userId, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("assetRepo.ListAssetCategory: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	return pkg.NewResponse(http.StatusOK, "success", data, nil)
}
