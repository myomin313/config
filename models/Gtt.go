package models

import "time"

type GameTypeTranslation struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;unsigned"`

	MasterGameTypeID uint64 `gorm:"not null;index;unsigned"` // FK to MasterGameType

	LanguageCode string `gorm:"size:10;not null;index"` // en-US, zh-CN, km-KH
	Name         string `gorm:"size:255;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Optional back reference
	MasterGameType MasterGameType `gorm:"foreignKey:MasterGameTypeID"`
}

func init() {
	RegisterModel(&GameTypeTranslation{})
}