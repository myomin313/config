package models

import "time"

type GameRedLimit struct {
	ID                   uint              `gorm:"primaryKey"`
	GameType             uint              `json:"game_type" gorm:"column:game_type"`
	GameTypeCurrency     uint              `json:"game_type_currency" gorm:"column:game_type_currency"`
	GamePlayType         uint              `json:"game_play_type" gorm:"column:game_play_type"`
	GamePlayTypeCurrency uint              `json:"game_play_type_currency" gorm:"column:game_play_type_currency"`
	Status               string            `json:"status" gorm:"default:disabled"`
	Sort                 uint              `json:"sort"`
	CreatedByID          uint              `json:"created_by_id"`
	UpdatedByID          uint              `json:"updated_by_id"`
	Translations         []GameTranslation `gorm:"polymorphic:Table;"`
	RedLimitTable        []RedLimitTable   `gorm:"foreignKey:GameRedLimitID;references:ID"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	GameAutoExits        uint              `json:"no_auto_exits" gorm:"->;column:no_auto_exits"`
}

func init() {
	RegisterModel(&GameRedLimit{})
}

func (GameRedLimit) DropColumns() []string {
	return []string{
		"pc_icon",
	}
}
