package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

type YahooProvider interface {
	FetchPrice(ctx context.Context, ticker string) (decimal.Decimal, error)
}
