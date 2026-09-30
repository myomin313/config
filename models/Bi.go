package models

import (
	"time"
)

func init() {
	RegisterModel(&BannerImage{})
}

type BannerImage struct {
	ID           uint64  `gorm:"primaryKey"`
	BannerId     uint64  `gorm:"not null;index"` // indexed
	ImagePath    *string `gorm:"type:varchar(255);not null"`
	LanguageCode string  `gorm:"type:varchar(255);not null"`
	DefaultImage *bool   `gorm:"not null;default:false"` // fixed default syntax
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
