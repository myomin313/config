package models

import "time"

type Game struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement"`
	BettingTime    int    `gorm:"not null;default:0"`
	SettlementTime int    `gorm:"not null;default:0"`
	NoAutoExits    int    `gorm:"default:0"`
	GameTypeID uint64 `gorm:"not null;default:0;index"`
    Status string `gorm:"default:disable"`
	CreatedBy uint
	UpdatedBy uint
    GameDataID uint
	// Relations to User
	// CreatedByUser *User `gorm:"foreignKey:CreatedBy;references:ID"`
	// UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;references:ID"`

	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`

	//Translations []GameTranslation `gorm:"foreignKey:GameID;references:ID"`
	Translations []GameTranslation `gorm:"polymorphic:Table;"`
	GameTypeData []GameTypeData   `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`

	//GameTypes []GameType `gorm:"many2many:game_game_types;constraint:OnDelete:CASCADE;"`
	GameTypes    []GameType     `gorm:"many2many:game_game_types;"`
}

// type GameGameType struct {
// 	GameID     uint64
// 	GameTypeID uint64
// }

func init() {
	RegisterModel(&Game{})
	//RegisterModel(&GameGameType{})
}
