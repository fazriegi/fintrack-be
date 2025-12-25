package usecase

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/database"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/repository"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/sirupsen/logrus"
)

type AssetUsecase interface {
	ListAssetCategory(userId uint) (resp pkg.Response)
	SubmitAsset(param entity.SubmitAssetRequest) (resp pkg.Response)
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

func (u *assetUsecase) SubmitAsset(param entity.SubmitAssetRequest) (resp pkg.Response) {
	var (
		err error
		db  = database.Get()
	)

	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error start transaction: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}
	defer tx.Rollback()

	dataInsert := entity.Asset{
		Name:          param.Name,
		CategoryId:    param.CategoryId,
		UserId:        param.UserId,
		Amount:        param.Amount,
		PurchasePrice: param.PurchasePrice,
		Status:        param.Status,
	}

	err = u.assetRepo.Insert(dataInsert, tx)
	if err != nil {
		u.log.Errorf("assetRepo.Insert: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	return pkg.NewResponse(http.StatusCreated, "success", param, nil)
}
