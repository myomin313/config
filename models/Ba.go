package models

import (
	enums "config-service/Enums"
	"time"
)

type Banner struct {
	ID              uint64                `gorm:"primaryKey"`
	ImageName       string                `gorm:"type:varchar(255);not null"`
	DisplayPlatform enums.DisplayPlatform `gorm:"not null"`
	DisplayDuration enums.DisplayDuration `gorm:"not null"` // "schedule" or "permanent"
	StartDate       *time.Time
	EndDate         *time.Time
	Status          enums.Status `gorm:"not null;default:1"`
	CreatedByID     uint64       `json:"created_by_id"`
	UpdatedByID     uint64       `json:"updated_by_id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	// Relation
	BannerImages []BannerImage `gorm:"foreignKey:BannerId;constraint:OnDelete:CASCADE"`
}

func init() {
	RegisterModel(&Banner{})
}
