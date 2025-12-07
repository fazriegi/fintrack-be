package entity

type User struct {
	ID       uint   `db:"id" json:"id"`
	Name     string `db:"name" json:"name" validate:"required,min=2,max=100"`
	Email    string `db:"email" json:"email" validate:"email"`
	Username string `db:"username" json:"username" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

type UserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
