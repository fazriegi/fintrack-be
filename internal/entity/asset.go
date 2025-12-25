package entity

type (
	AssetCategory struct {
		ID   uint   `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
)
