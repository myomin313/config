package models

import "time"

type Maintenance struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	ApplicationMaintenance bool `gorm:"not null;default:false" json:"application_maintenance"`
	GameMaintenance        bool `gorm:"not null;default:false" json:"game_maintenance"`

	CreatedByID *uint `gorm:"index" json:"created_by_id,omitempty"`
	UpdatedByID *uint `gorm:"index" json:"updated_by_id,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func init() {
	RegisterModel(&Maintenance{})
}

// func (Maintenance) DropColumns() []string {
// 	return []string{}
// }
