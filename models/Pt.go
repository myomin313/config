package models

import (
	"time"

	"gorm.io/datatypes"
)

type Platform struct {
	ID           uint64                      `gorm:"primaryKey;autoIncrement"`
	MerchandCode string                      `gorm:"size:255"`
	Name         string                      `gorm:"size:255;not null"`
	Email        string                      `gorm:"size:255;unique"`
	Status       string                      `gorm:"size:16;not null;default:enable"`
	IPAllowlist  datatypes.JSONSlice[string] `gorm:"type:jsonb"`
	CreatedAt    time.Time                   `gorm:"autoCreateTime"`
	UpdatedAt    time.Time                   `gorm:"autoUpdateTime"`

	// One-to-one relation
	APIKey             *PlatformApiKey    `gorm:"foreignKey:PlatformID;references:ID"`
	MerchantCurrencies []MerchantCurrency `gorm:"foreignKey:PlatformID"`
}

func init() {
	RegisterModel(&Platform{})
}
