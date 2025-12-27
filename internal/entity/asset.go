package entity

import (
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/shopspring/decimal"
)

type (
	AssetCategory struct {
		ID   uint   `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}

	Asset struct {
		Id            uint            `db:"id" json:"id"`
		Name          string          `db:"name" json:"name"`
		CategoryId    uint            `db:"category_id" json:"-"`
		UserId        uint            `db:"user_id" json:"-"`
		Amount        decimal.Decimal `db:"amount" json:"amount"`
		PurchasePrice decimal.Decimal `db:"purchase_price" json:"purchase_price"`
		Status        string          `db:"status" json:"status"`
	}

	AssetResponse struct {
		Id                 uint            `db:"id" json:"id"`
		Name               string          `db:"name" json:"name"`
		CategoryId         uint            `db:"category_id" json:"-"`
		Category           string          `db:"category" json:"category"`
		UserId             uint            `db:"user_id" json:"-"`
		Amount             decimal.Decimal `db:"amount" json:"amount"`
		PurchasePrice      decimal.Decimal `db:"purchase_price" json:"purchase_price"`
		TotalPurchasePrice decimal.Decimal `db:"total_purchase_price" json:"total_purchase_price"`
		Status             string          `db:"status" json:"status"`
	}

	SubmitAssetRequest struct {
		UserId        uint            `json:"-"`
		Name          string          `json:"name" validate:"required,min=1,max=100"`
		CategoryId    uint            `json:"category_id" validate:"required"`
		Amount        decimal.Decimal `json:"amount" validate:"required,decimal_gt_zero"`
		PurchasePrice decimal.Decimal `json:"purchase_price" validate:"required,decimal_gt_zero"`
		Status        string          `json:"status" validate:"required,oneof=active inactive sold"`
	}

	ListAssetRequest struct {
		pkg.PaginationRequest
		UserId   uint
		Name     string `query:"name"`
		Category string `query:"category"`
	}
)
