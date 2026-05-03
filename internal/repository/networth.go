package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type networthRepository struct{}

type NetworthRepository interface {
	Calculate(ctx context.Context, db *sqlx.DB) error
}

func NewNetworthRepository() NetworthRepository {
	return &networthRepository{}
}

func (r *networthRepository) Calculate(ctx context.Context, db *sqlx.DB) error {
	query := `
		WITH AssetSummary AS (
			SELECT user_id, COALESCE(SUM(current_value), 0) AS total_assets
			FROM assets 
			WHERE is_active = TRUE 
			GROUP BY user_id
		),
		LiabilitySummary AS (
			SELECT user_id, COALESCE(SUM(remaining_balance), 0) AS total_liabilities
			FROM liabilities 
			WHERE remaining_balance > 0 
			GROUP BY user_id
		)
		INSERT INTO net_worth_histories (user_id, total_assets, total_liabilities, recorded_date)
		SELECT 
			U.id, 
			COALESCE(A.total_assets, 0), 
			COALESCE(L.total_liabilities, 0),
			CURRENT_DATE
		FROM users U
		LEFT JOIN AssetSummary A ON U.id = A.user_id
		LEFT JOIN LiabilitySummary L ON U.id = L.user_id
		ON CONFLICT (user_id, recorded_date) 
		DO UPDATE SET 
			total_assets = EXCLUDED.total_assets,
			total_liabilities = EXCLUDED.total_liabilities,
			updated_at = NOW();`

	_, err := db.ExecContext(ctx, query)

	return err
}
