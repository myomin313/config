package requests

// GameTableCreateRequest represents the request body for creating a Platform
// GameTableCreateRequest represents the request body for creating a Game Table

type GameTableCreateRequest struct {
	GameTableName   map[string]string `json:"game_table_name" validate:"required,require_all_game_locales" example:"{\"en\":\"Dragon Tiger\",\"km\":\"ដ្រាកន៍ ថាយហ្គ័\"}"`
	AffiliateOffice uint              `json:"game_hall" validate:"required,notblank" example:"1"`
	GameType        uint              `json:"game" validate:"required,notblank" example:"2"`
	GameTypeID      uint              `json:"game_type_id" validate:"required,gt=0" example:"1"`
	BettingTime     uint              `json:"betting_time" validate:"required,gt=0" example:"0"`
	SettlementTime  uint              `json:"settlement_time" validate:"required,gt=0" example:"0"`
	//TableNumber     string            `json:"table_number" validate:"required,notblank,nohyphen" example:"T_01"`
	TableNumber string `json:"table_number" validate:"required,table_number" example:"T_01"`
	Status      string `json:"status" validate:"required,notblank,oneof=active  disabled  maintenance" example:"active"`
	//TableResultType string                       `json:"table_result_type" validate:"required,notblank" example:"animal_face"`
	MobileIcon  string                       `json:"mobile_icon" validate:"required,notblank" example:"https://example.com/mobile.png"`
	PCIcon      string                       `json:"pc_icon" validate:"required,notblank" example:"https://example.com/pc.png"`
	Sort        uint                         `json:"sort" validate:"required,notblank,gt=0" example:"0"`
	Mode        string                       `json:"mode" validate:"required,notblank" example:"dealer"`
	Topic       string                       `json:"topic"`
	Resolutions []GameTableResolutionRequest `json:"resolutions"`
}

// GameTableUpdateRequest represents the request body for updating a GameTable
type GameTableUpdateRequest struct {
	GameTableName *map[string]string `json:"game_table_name" validate:"omitempty,required,require_all_game_locales" example:"{\"en\":\"Dragon Tiger\",\"km\":\"ដ្រាកន៍ ថាយហ្គ័\"}"`

	AffiliateOffice *uint `json:"game_hall" validate:"omitempty,required,notblank" example:"1"`
	GameType        *uint `json:"game" validate:"omitempty,required,notblank" example:"2"`
	GameTypeID      *uint `json:"game_type_id" validate:"omitempty,gt=0" example:"1"`
	BettingTime     *uint `json:"betting_time" validate:"omitempty,required,gt=0" example:"0"`
	SettlementTime  *uint `json:"settlement_time" validate:"omitempty,required,gt=0" example:"0"`
	//TableNumber     *string `json:"table_number" validate:"omitempty,required,notblank,nohyphen" example:"T_01"`
	TableNumber *string `json:"table_number" validate:"required,table_number" example:"T_01"`
	Status      *string `json:"status" validate:"omitempty,required,notblank,oneof=active disabled maintenance" example:"disabled"`
	//TableResultType *string                             `json:"table_result_type" validate:"omitempty,required,notblank" example:"animal_face"`
	MobileIcon  *string                             `json:"mobile_icon" validate:"omitempty,required,notblank" example:"https://example.com/mobile.png"`
	PCIcon      *string                             `json:"pc_icon" validate:"omitempty,required,notblank" example:"https://example.com/pc.png"`
	Sort        *uint                               `json:"sort" validate:"omitempty,required,notblank,gt=0" example:"0"`
	Mode        *string                             `json:"mode" validate:"required,notblank" example:"dealer"`
	Topic       *string                             `json:"topic"`
	Resolutions *[]GameTableResolutionUpdateRequest `json:"resolutions"`
}

type GameTableSearchRequest struct {
	TableName *string `query:"table_name" validate:"omitempty,notblank" example:"dragon-tiger"`
	GameType  *uint   `query:"game" validate:"omitempty,gt=0" example:"1"`
	Status    *string `query:"status" validate:"omitempty,oneof=active  disabled  maintenance"`
}

type GameTableMassStatusUpdateRequest struct {
	TableIDs []uint `json:"table_ids" validate:"required,min=1" example:"1,2,3"`
	Status   string `json:"status" validate:"required,oneof=active  disabled  maintenance" example:"active"`
}

type GameTableFESearchRequest struct {
	AffiliateOffice *uint `query:"game_hall" validate:"omitempty,required,notblank" example:"1"`
	Game            *uint `query:"game" validate:"required,gt=0" example:"1"`
}
type GameTableResolutionRequest struct {
	Resolution int    `json:"resolution" validate:"required"`
	RoomName   string `json:"room_name" validate:"required"`
}
type GameTableResolutionUpdateRequest struct {
	Resolution int    `json:"resolution" validate:"required,gt=0"`
	RoomName   string `json:"room_name" validate:"required,notblank"`
}
type GameTableBasicUpdateRequest struct {
	GameTableName   *map[string]string `json:"game_table_name"`
	AffiliateOffice *uint              `json:"game_hall"`
	GameType        *uint              `json:"game"`
	GameTypeID      *uint              `json:"game_type_id"`
	BettingTime     *uint              `json:"betting_time"`
	SettlementTime  *uint              `json:"settlement_time"`
	TableNumber     *string            `json:"table_number"`
	Status          *string            `json:"status"`
	//TableResultType *string            `json:"table_result_type"`
	MobileIcon *string `json:"mobile_icon"`
	PCIcon     *string `json:"pc_icon"`
	Sort       *uint   `json:"sort"`
	Mode       *string `json:"mode"`
	Topic      *string `json:"topic"`
}
