package models

import "time"

type MerchantCurrency struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantCode string    `gorm:"type:varchar(50);not null;index" json:"merchant_code"`
	Currency     string    `gorm:"type:varchar(50);not null" json:"currency"`
	PlatformID   uint      `gorm:"not null;index" json:"platform_id"`
	Status       string    `gorm:"size:16;not null;default:enable"`
	Platform     Platform  `gorm:"foreignKey:PlatformID;references:ID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func init() {
	RegisterModel(&MerchantCurrency{})
}
