package models

import "time"

type GameTable struct {
	ID              uint   `gorm:"primaryKey"`
	BettingTime     uint   `json:"betting_time"`
	SettlementTime  uint   `json:"settlement_time"`
	AffiliateOffice uint   `json:"affiliate_office"`                        //game hall
	GameType        uint   `json:"game_type" gorm:"column:game_type"`       //gameId
	GameTypeID      uint   `json:"game_type_id" gorm:"column:game_type_id"` //gamePlayTypeId
	TableNumber     string `json:"table_number"`
	Status          string `json:"status" gorm:"default:disabled"`
	MobileIcon      string `json:"mobile_icon"`
	//TableResultType   string                `json:"table_result_type" gorm:"type:varchar(50);default:'animal_face'"`
	Sort              uint                  `json:"sort"`
	PCIcon            string                `json:"pc_icon"`
	CreatedByID       uint                  `json:"created_by_id"`
	UpdatedByID       uint                  `json:"updated_by_id"`
	Translations      []GameTranslation     `gorm:"polymorphic:Table;"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	Mode              string                `json:"mode"`
	Topic             string                `json:"topic"`
	KafkaOnlineStatus *uint                 `json:"kafka_online_status"`
	Resolutions       []GameTableResolution `gorm:"foreignKey:GameTableID;constraint:OnDelete:CASCADE" json:"resolutions"`
}

func init() {
	RegisterModel(&GameTable{})
}

func (GameTable) DropColumns() []string {
	return []string{
		"video_type",
		"room_name",
		"video_url",
		"idle_households",
		"table_color",
	}
}
