package entity

import "github.com/shopspring/decimal"

type (
	AssetCategory struct {
		ID   uint   `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}

	Asset struct {
		Name          string          `db:"name"`
		CategoryId    uint            `db:"category_id"`
		UserId        uint            `db:"user_id"`
		Amount        decimal.Decimal `db:"amount"`
		PurchasePrice decimal.Decimal `db:"purchase_price"`
		Status        string          `db:"status"`
	}

	SubmitAssetRequest struct {
		UserId        uint            `json:"-"`
		Name          string          `json:"name" validate:"required,min=1,max=100"`
		CategoryId    uint            `json:"category_id" validate:"required"`
		Amount        decimal.Decimal `json:"amount" validate:"required"`
		PurchasePrice decimal.Decimal `json:"purchase_price" validate:"required"`
		Status        string          `json:"status" validate:"required,oneof=active inactive sold"`
	}
)
