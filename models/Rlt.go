package models

import "time"

type RedLimitTable struct {
	ID              uint   `gorm:"primaryKey"`
	GameRedLimitID  uint64 `json:"red_limit_id" gorm:"column:red_limit_id;not null;index;unsigned"`
	GameRedLimitKey string `json:"red_limit_key" gorm:"column:red_limit_key;not null;"`
	// MinLimit        uint      `json:"min_limit" gorm:"column:min_limit"`
	// MaxLimit        uint      `json:"max_limit" gorm:"column:max_limit"`
	MinLimit  float64   `json:"min_limit" gorm:"column:min_limit;type:decimal(10,2)"`
	MaxLimit  float64   `json:"max_limit" gorm:"column:max_limit;type:decimal(10,2)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	RegisterModel(&RedLimitTable{})
}
