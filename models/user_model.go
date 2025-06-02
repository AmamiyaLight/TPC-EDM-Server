package models

import "TPC-EDM-Server/models/enum"

type UserModel struct {
	Model
	Username string        `gorm:"size:32" json:"username"`
	Password string        `gorm:"size:64" json:"-"`
	Role     enum.RoleType `json:"role"`
}
