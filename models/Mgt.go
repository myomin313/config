package models

import "time"

type MasterGameType struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`

	Odds      int       `gorm:"not null;default:1"`
	Image     string    `gorm:"size:255"`
	KeyName   string    `gorm:"size:255"`

	// NameEn    string    `gorm:"size:255;not null"`
	// NameZh    string    `gorm:"size:255;not null"`
	// NameKm    string    `gorm:"size:255;not null"`
    GameTypeID uint64
	GameType   GameType `gorm:"foreignKey:GameTypeID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	GameTypeTranslation []GameTypeTranslation `gorm:"foreignKey:MasterGameTypeID"`
}
func init() {
	RegisterModel(&MasterGameType{})
}