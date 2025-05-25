package models

import "time"

type Model struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id" binding:"required"`
}

type RemoveRequest struct {
	IDList []uint `json:"id_list" binding:"required"`
}

type OptionsResponse[T any] struct {
	Label string `json:"label"`
	Value T      `json:"value"`
}
