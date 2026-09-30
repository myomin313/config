package models

import "time"

type GameTranslation struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;unsigned"`
	TableID       uint64    `json:"table_id" gorm:"column:table_id;not null;index;unsigned"`
	LanguageCode string    `gorm:"size:191;not null;index"`
	Content      string    `gorm:"type:text;not null"`
	Type         string    `gorm:"size:191;not null;default:api"`
	LangName     string    `gorm:"size:191;not null;default:game_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TableType string `gorm:"size:191;not null"`
}

func init() {
	RegisterModel(&GameTranslation{})
}
