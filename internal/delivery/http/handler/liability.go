package handler

import (
	"log"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/delivery/http/middleware"
	"github.com/fazriegi/fintrack-be/internal/usecase"
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
}

func (h *LiabilityHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	response := h.usecase.ListCategory(r.Context())
	response.HTTP(w)
}
