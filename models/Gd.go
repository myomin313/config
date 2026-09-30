package models

import "time"

type GameData struct {
	ID              uint              `gorm:"primaryKey"`
	NameCh     string            `json:"Name_ch"`
	NameEn          string            `json:"Name_en"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func init() {
	RegisterModel(&GameData{})
}
