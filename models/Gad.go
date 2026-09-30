package models

import "time"

type GameTypeData struct {
	ID               uint64  `gorm:"primaryKey;autoIncrement;unsigned"`
	MasterGameTypeID uint64  `gorm:"not null;index;unsigned"`
	GameID           uint64  `gorm:"not null;index;unsigned"`
	Image            string  `gorm:"size:255"`
	NameKm           string  `gorm:"size:255;not null"`
	Odds             float64 `gorm:"type:decimal(10,3);not null;default:1.000"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Game           *Game           `gorm:"foreignKey:GameID;references:ID"`
	MasterGameType *MasterGameType `gorm:"foreignKey:MasterGameTypeID;references:ID"`
}

func init() {
	RegisterModel(&GameTypeData{})
}
