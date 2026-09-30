package models

import "time"

type GameTableResolution struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GameTableID uint      `gorm:"not null;index" json:"game_table_id"`
	Resolution  int       `gorm:"not null" json:"resolution"`
	RoomName    string    `gorm:"type:varchar(255);not null" json:"room_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func init() {
	RegisterModel(&GameTableResolution{})
}
