package pkg

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/database"
)

func QueryWithPagination(dataset *goqu.SelectDataset, req PaginationRequest) *goqu.SelectDataset {
	if req.Sort != nil && *req.Sort != "" {
		parts := strings.Fields(*req.Sort)
		if len(parts) > 0 {
			field := parts[0]
			direction := "ASC"
			if len(parts) >= 2 && strings.ToUpper(parts[1]) == "DESC" {
				direction = "DESC"
			}

			col := goqu.I(field)
			if direction == "DESC" {
				dataset = dataset.Order(col.Desc())
			} else {
				dataset = dataset.Order(col.Asc())
			}
		}
	}

	if req.Page != nil && req.Limit != nil && *req.Page > 0 && *req.Limit > 0 {
		offset := (*req.Page - 1) * *req.Limit
		dataset = dataset.Limit(*req.Limit).Offset(offset)
	}

	return dataset
}

func GetDialect() goqu.DialectWrapper {
	return goqu.Dialect(database.GetDriver())
}
