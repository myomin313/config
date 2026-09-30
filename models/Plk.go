package models

import "time"

type PlatformApiKey struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	PlatformID uint64    `gorm:"not null;uniqueIndex;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status     string    `gorm:"size:16;not null;default:disabled"`
	Key        string    `gorm:"size:100;unique;not null"`
	Secret     string    `gorm:"type:text;unique;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Platform *Platform `gorm:"foreignKey:PlatformID;references:ID"`
}

func init() {
	RegisterModel(&PlatformApiKey{})
}
