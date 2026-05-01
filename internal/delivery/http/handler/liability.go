package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/delivery/http/middleware"
	"github.com/fazriegi/fintrack-be/internal/domain"
	"github.com/fazriegi/fintrack-be/internal/usecase"
	"github.com/fazriegi/fintrack-be/pkg"
	"github.com/fazriegi/fintrack-be/pkg/constant"
	"github.com/fazriegi/fintrack-be/pkg/validator"
	"github.com/google/uuid"
)

type LiabilityHandler struct {
	usecase usecase.LiabilityUsecase
	logger  *log.Logger
}

func NewLiabilityHandler(mux *http.ServeMux, uc usecase.LiabilityUsecase, logger *log.Logger) {
	handler := &LiabilityHandler{
		usecase: uc,
		logger:  logger,
	}

	mux.Handle("GET /v1/liabilities/categories", middleware.MiddlewareAuth(http.HandlerFunc(handler.ListCategory)))
	mux.Handle("POST /v1/liabilities", middleware.MiddlewareAuth(http.HandlerFunc(handler.Create)))
	mux.Handle("GET /v1/liabilities", middleware.MiddlewareAuth(http.HandlerFunc(handler.List)))
	mux.Handle("GET /v1/liabilities/{id}", middleware.MiddlewareAuth(http.HandlerFunc(handler.GetByID)))
	mux.Handle("PUT /v1/liabilities/{id}", middleware.MiddlewareAuth(http.HandlerFunc(handler.Update)))
	mux.Handle("DELETE /v1/liabilities/{id}", middleware.MiddlewareAuth(http.HandlerFunc(handler.Delete)))
}

func (h *LiabilityHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	response := h.usecase.ListCategory(r.Context())
	response.HTTP(w)
}

func (h *LiabilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateLiability

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.NewResponse(http.StatusBadRequest, constant.ErrInvalidJson, nil, nil).HTTP(w)
		return
	}

	// validation
	validationErr := validator.ValidateRequest(&req)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		pkg.NewResponse(http.StatusUnprocessableEntity, constant.ErrValidation, errResponse, nil).HTTP(w)
		return
	}

	response := h.usecase.Create(r.Context(), &req)
	response.HTTP(w)
}

func (h *LiabilityHandler) List(w http.ResponseWriter, r *http.Request) {
	var req domain.ListLiabilityRequest

	if err := pkg.ParseQueryParam(r, &req); err != nil {
		h.logger.Printf("[ERROR] parsing query params: %s", err.Error())
		pkg.NewResponse(http.StatusBadRequest, constant.ErrParseQueryParam, nil, nil).HTTP(w)
		return
	}

	response := h.usecase.List(r.Context(), &req)

	response.HTTP(w)
}

func (h *LiabilityHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Printf("[ERROR] uuid.Parse - invalid UUID format: %s", err.Error())
		pkg.NewResponse(http.StatusBadRequest, constant.ErrInvalidParam, nil, nil).HTTP(w)
		return
	}

	h.usecase.GetByID(r.Context(), parsedID).HTTP(w)
}

func (h *LiabilityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Printf("[ERROR] uuid.Parse - invalid UUID format: %s", err.Error())
		pkg.NewResponse(http.StatusBadRequest, constant.ErrInvalidParam, nil, nil).HTTP(w)
		return
	}

	var req domain.CreateLiability

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.NewResponse(http.StatusBadRequest, constant.ErrInvalidJson, nil, nil).HTTP(w)
		return
	}

	// validation
	validationErr := validator.ValidateRequest(&req)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		pkg.NewResponse(http.StatusUnprocessableEntity, constant.ErrValidation, errResponse, nil).HTTP(w)
		return
	}

	req.ID = parsedID

	h.usecase.Update(r.Context(), &req).HTTP(w)
}

func (h *LiabilityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Printf("[ERROR] uuid.Parse - invalid UUID format: %s", err.Error())
		pkg.NewResponse(http.StatusBadRequest, constant.ErrInvalidParam, nil, nil).HTTP(w)
		return
	}

	h.usecase.Delete(r.Context(), parsedID).HTTP(w)
}
