package models

import "time"

type Announcement struct {
	ID                         uint              `gorm:"primaryKey"`
	AnnouncementName           string            `json:"announcement_name"`
	Status                     string            `json:"status" gorm:"default:disabled"`
	SortOrder                  uint              `json:"sort_order"`
	BroadcastAudience          string            `json:"broadcast_audience"`
	CreatedBy                  uint              `json:"created_by"`
	UpdatedBy                  uint              `json:"updated_by"`
	CreatedAt                  time.Time         `json:"created_at"`
	UpdatedAt                  time.Time         `json:"updated_at"`
	AnnouncementDetail         []AnnouncementDetail `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE"`
}
type AnnouncementDetail struct {
	ID                         uint              `gorm:"primaryKey"`
	AnnouncementID             uint              `json:"announcement_id"`
	LanguageCode               string            `json:"language_code"`
	Title                      string            `json:"title"`
	Content                    string            `gorm:"type:text"`
	CreatedAt                  time.Time         `json:"created_at"`
	UpdatedAt                  time.Time         `json:"updated_at"`
}

func init() {
	RegisterModel(&Announcement{})
	RegisterModel(&AnnouncementDetail{})
}
