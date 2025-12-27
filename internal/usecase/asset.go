package usecase

import (
	"database/sql"
	"errors"
	"math"
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
	ListAsset(param entity.ListAssetRequest) (resp pkg.Response)
	SubmitAsset(param entity.SubmitAssetRequest) (resp pkg.Response)
	GetById(param entity.GetAssetByIdRequest) (resp pkg.Response)
	Update(param entity.UpdateAssetRequest) (resp pkg.Response)
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

func (u *assetUsecase) ListAsset(param entity.ListAssetRequest) (resp pkg.Response) {
	var (
		err error
		db  = database.Get()
	)

	data, total, err := u.assetRepo.ListAsset(param, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("assetRepo.ListAsset: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	var paginationMeta pkg.PaginationMeta
	if param.Limit != nil && *param.Limit > 0 {
		limit := int(*param.Limit)
		page := 1

		if param.Page != nil && *param.Page > 0 {
			page = int(*param.Page)
		}

		totalPages := int(math.Ceil(float64(total) / float64(limit)))

		if totalPages > 0 && page > totalPages {
			page = totalPages
		}

		paginationMeta = pkg.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		}
	}

	return pkg.NewResponse(http.StatusOK, "success", data, &paginationMeta)
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

func (u *assetUsecase) GetById(param entity.GetAssetByIdRequest) (resp pkg.Response) {
	var (
		err error
		db  = database.Get()
	)

	data, err := u.assetRepo.GetById(param.Id, param.UserId, false, db, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.NewResponse(http.StatusNotFound, pkg.ErrNotFound.Error(), nil, nil)
		}

		u.log.Errorf("assetRepo.GetById: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	return pkg.NewResponse(http.StatusOK, "success", data, nil)
}

func (u *assetUsecase) Update(param entity.UpdateAssetRequest) (resp pkg.Response) {
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

	data, err := u.assetRepo.GetById(param.Id, param.UserId, true, nil, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.NewResponse(http.StatusNotFound, pkg.ErrNotFound.Error(), nil, nil)
		}

		u.log.Errorf("assetRepo.GetById: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	dataUpdate := entity.Asset{
		Id:            param.Id,
		Name:          param.Name,
		CategoryId:    param.CategoryId,
		UserId:        param.UserId,
		Amount:        param.Amount,
		PurchasePrice: param.PurchasePrice,
		Status:        param.Status,
	}

	err = u.assetRepo.Update(dataUpdate, param.Id, param.UserId, tx)
	if err != nil {
		u.log.Errorf("assetRepo.Update: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	data = entity.AssetResponse{
		Id:            param.Id,
		Name:          param.Name,
		CategoryId:    param.CategoryId,
		UserId:        param.UserId,
		Amount:        param.Amount,
		PurchasePrice: param.PurchasePrice,
		Status:        param.Status,
	}

	return pkg.NewResponse(http.StatusOK, "success", data, nil)
}
