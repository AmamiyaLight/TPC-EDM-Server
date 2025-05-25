package models

type UserModel struct {
	Model
	Username string `gorm:"size:32" json:"username"`
	Password string `gorm:"size:64" json:"-"`
	Role     int8   `json:"role"`
}
