package usecase

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/domain"
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
	Create(ctx context.Context, req *domain.CreateLiability) (resp pkg.Response)
	List(ctx context.Context, req *domain.ListLiabilityRequest) (resp pkg.Response)
	GetByID(ctx context.Context, id uuid.UUID) (resp pkg.Response)
	Update(ctx context.Context, req *domain.CreateLiability) (resp pkg.Response)
	Delete(ctx context.Context, id uuid.UUID) (resp pkg.Response)
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

func (u *liabilityUsecase) Create(ctx context.Context, req *domain.CreateLiability) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)
	req.UserId = userId

	var detailsDB any
	if req.Details != nil {
		b, err := json.Marshal(req.Details)
		if err != nil {
			u.log.Printf("[ERROR] json.Marshal req.Details: %s", err.Error())
			return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
		}
		detailsDB = b
	}

	liabilityDB := &domain.LiabilityDB{
		UserId:           userId,
		CategoryID:       req.CategoryID,
		Name:             req.Name,
		PrincipalAmount:  *req.PrincipalAmount,
		RemainingBalance: *req.RemainingBalance,
		Details:          detailsDB,
	}

	err := u.repo.Insert(ctx, liabilityDB, u.db)
	if err != nil {
		u.log.Printf("[ERROR] repo.Insert: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
	}

	return pkg.NewResponse(http.StatusCreated, "Success", nil, nil)
}

func (u *liabilityUsecase) List(ctx context.Context, req *domain.ListLiabilityRequest) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)
	req.UserId = userId

	liabilities, total, err := u.repo.List(ctx, req, u.db)
	if err != nil {
		u.log.Printf("[ERROR] repo.List: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
	}

	dataResponse := make([]domain.ListLiabilityResponse, 0)

	if liabilities != nil {
		for _, liability := range *liabilities {
			dataResponse = append(dataResponse, domain.ListLiabilityResponse{
				ID:               liability.ID,
				Category:         liability.Category,
				Name:             liability.Name,
				RemainingBalance: liability.RemainingBalance,
			})
		}
	}

	var paginationMeta pkg.PaginationMeta
	if req.Limit != nil && *req.Limit > 0 {
		limit := int(*req.Limit)
		page := 1

		if req.Page != nil && *req.Page > 0 {
			page = int(*req.Page)
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

	return pkg.NewResponse(http.StatusOK, "Success", dataResponse, &paginationMeta)
}

func (u *liabilityUsecase) GetByID(ctx context.Context, id uuid.UUID) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)

	liability, err := u.repo.GetByID(ctx, id, userId, u.db)
	if err != nil {
		if err.Error() != constant.ErrNotFound {
			u.log.Printf("[ERROR] repo.GetByID: %s", err.Error())
			return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
		}

		return pkg.NewResponse(http.StatusNotFound, constant.ErrNotFound, nil, nil)
	}

	dataResponse := domain.GetLiabilityByIDResponse{
		ID:               liability.ID,
		CategoryID:       liability.CategoryID,
		Category:         liability.Category,
		CategoryType:     liability.CategoryType,
		Name:             liability.Name,
		PrincipalAmount:  liability.PrincipalAmount,
		RemainingBalance: liability.RemainingBalance,
	}

	if liability.Details != nil {
		var detailsBytes []byte
		switch v := liability.Details.(type) {
		case []byte:
			detailsBytes = v
		case string:
			detailsBytes = []byte(v)
		}

		if len(detailsBytes) > 0 {
			var parsedDetails any
			if err := json.Unmarshal(detailsBytes, &parsedDetails); err == nil {
				dataResponse.Details = parsedDetails
			}
		}
	}

	return pkg.NewResponse(http.StatusOK, "Success", dataResponse, nil)
}

func (u *liabilityUsecase) Update(ctx context.Context, req *domain.CreateLiability) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)
	req.UserId = userId

	var detailsDB any
	if req.Details != nil {
		b, err := json.Marshal(req.Details)
		if err != nil {
			u.log.Printf("[ERROR] json.Marshal req.Details: %s", err.Error())
			return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
		}
		detailsDB = b
	}

	liabilityDB := &domain.LiabilityDB{
		ID:               req.ID,
		UserId:           userId,
		CategoryID:       req.CategoryID,
		Name:             req.Name,
		PrincipalAmount:  *req.PrincipalAmount,
		RemainingBalance: *req.RemainingBalance,
		Details:          detailsDB,
	}

	err := u.repo.Update(ctx, liabilityDB, u.db)
	if err != nil {
		u.log.Printf("[ERROR] repo.Update: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
	}

	return pkg.NewResponse(http.StatusOK, "Success", nil, nil)
}

func (u *liabilityUsecase) Delete(ctx context.Context, id uuid.UUID) (resp pkg.Response) {
	userId := ctx.Value("user_id").(uuid.UUID)

	err := u.repo.Delete(ctx, id, userId, u.db)
	if err != nil {
		u.log.Printf("[ERROR] repo.Delete: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, constant.ErrServer, nil, nil)
	}

	return pkg.NewResponse(http.StatusOK, "Success", nil, nil)
}
