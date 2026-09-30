package requests

type GameRequest struct {
	Name           string `json:"name" validate:"required,min=3"`
	BettingTime    int    `json:"betting_time" validate:"required,min=1"`
	SettlementTime int    `json:"settlement_time" validate:"required,min=1"`
	NoAutoExits    int    `json:"no_auto_exits"`
	CreatedBy      uint64 `json:"created_by" validate:"required"`
	UpdatedBy      uint64 `json:"updated_by" validate:"required"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type GameCreateRequest struct {
	BettingTime    int               `json:"betting_time" validate:"gte=1"`
	SettlementTime int               `json:"settlement_time" validate:"gte=1"`
	NoAutoExits    int               `json:"no_auto_exits"`
	GameDataID     uint              `json:"game_data_id" validate:"required"`
	GameTypeIDs    []uint64          `json:"game_type_ids" validate:"required,dive,gt=0"`
	Status         string            `json:"status" validate:"required"`
	CreatedBy      uint              `json:"created_by"`
	UpdatedBy      uint              `json:"updated_by"`
	Names          map[string]string `json:"names" validate:"required,require_all_game_locales"` // en, zh-CN, km

	GameTypeData []GameTypeDataRequest `json:"game_type_data" validate:"required,dive"`
}
type GameUpdateRequest struct {
	BettingTime    int    `json:"betting_time" validate:"gte=1"`
	SettlementTime int    `json:"settlement_time" validate:"gte=1"`
	NoAutoExits    int    `json:"no_auto_exits"`
	Status         string `json:"status" validate:"required"`
	//GameTypeID     uint64 `json:"game_type_id" validate:"required"`
	GameDataID  uint     `json:"game_data_id" validate:"required"`
	GameTypeIDs []uint64 `json:"game_type_ids" validate:"required,dive,gt=0"`

	UpdatedBy uint `json:"updated_by"`

	Names map[string]string `json:"names" validate:"required,require_all_game_locales"` // en, zh-CN, km

	GameTypeData []GameTypeDataRequest `json:"game_type_data" validate:"required,dive"`
}

type GameTypeRequest struct {
	Odds  float64 `json:"odds" validate:"required"`
	Image string  `json:"image" validate:"required"`

	Names map[string]string `json:"names" validate:"required"` // en, zh-CN, km
}

type GameTypeDataRequest struct {
	MasterGameTypeID uint64  `json:"master_game_type_id" validate:"required"`
	Odds             float64 `json:"odds" validate:"gt=0,max3dp"`
	NameKm           string  `json:"name_km" validate:"required"`
	Image            string  `json:"image"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
