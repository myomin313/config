package Resources

type GameTableResource struct {
	ID              uint   `json:"id"`
	ENName          string `json:"en_name"`
	CHName          string `json:"ch_name"`
	KHName          string `json:"kh_name"`
	AffiliateOffice uint   `json:"game_hall"`
	GameType        uint   `json:"game"`
	GameTypeID      uint   `json:"game_type_id"`
	BettingTime     uint   `json:"betting_time"`
	SettlementTime  uint   `json:"settlement_time"`
	TableNumber     string `json:"table_number"`
	//TableResultType   string                        `json:"table_result_type"`
	Status            string                        `json:"status"`
	MobileIcon        string                        `json:"mobile_icon"`
	PCIcon            string                        `json:"pc_icon"`
	CreatedBy         string                        `json:"created_by"`
	Sort              uint                          `json:"sort"`
	UpdatedBy         string                        `json:"updated_by"`
	CreatedAt         string                        `json:"created_at"`
	UpdatedAt         string                        `json:"updated_at"`
	Resolutions       []GameTableResolutionResource `json:"resolutions"`
	Mode              string                        `json:"mode"`
	Topic             string                        `json:"topic"`
	KafkaOnlineStatus *uint                         `json:"kafka_online_status"`
}
type GameTableFEResource struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	AffiliateOffice uint   `json:"game_hall"`
	GameType        uint   `json:"game"`
	BettingTime     uint   `json:"betting_time"`
	SettlementTime  uint   `json:"settlement_time"`
	TableNumber     string `json:"table_number"`
	//TableResultType string `json:"table_result_type"`
	Status     string `json:"status"`
	MobileIcon string `json:"mobile_icon"`
	PCIcon     string `json:"pc_icon"`
	Sort       uint   `json:"sort"`
}
type GameTableLogResource struct {
	ENName string `json:"en_name"`
	CHName string `json:"ch_name"`
	KHName string `json:"kh_name"`
	//AffiliateOffice uint                 `json:"game_hall"`
	//GameType        uint                 `json:"game"`
	//GameTypeID      uint                 `json:"game_type_id"`
	BettingTime    uint   `json:"betting_time"`
	SettlementTime uint   `json:"settlement_time"`
	TableNumber    string `json:"table_number"`
	//TableResultType   string                           `json:"table_result_type"`
	Status            string                           `json:"status"`
	MobileIcon        string                           `json:"mobile_icon"`
	PCIcon            string                           `json:"pc_icon"`
	CreatedBy         string                           `json:"created_by"`
	Sort              uint                             `json:"sort"`
	UpdatedBy         string                           `json:"updated_by"`
	CreatedAt         string                           `json:"created_at"`
	UpdatedAt         string                           `json:"updated_at"`
	GameHall          *GameHallLogResource             `json:"game_hall,omitempty"`
	Game              *GameLogResource                 `json:"game,omitempty"`
	GamePlayType      string                           `json:"game_play_type"`
	Mode              string                           `json:"mode"`
	Topic             string                           `json:"topic"`
	KafkaOnlineStatus *uint                            `json:"kafka_online_status"`
	Resolution        []GameTableResolutionLogResource `json:"resolution"`
}
type GameHallLogResource struct {
	ENName string `json:"en_name"`
	CHName string `json:"ch_name"`
	KHName string `json:"kh_name"`
}
type GameLogResource struct {
	ENName string `json:"en_name"`
	CHName string `json:"ch_name"`
	KHName string `json:"kh_name"`
}
type GameTableResolutionResource struct {
	ID          uint   `json:"id"`
	GameTableID uint   `json:"game_table_id"`
	Resolution  int    `json:"resolution"`
	RoomName    string `json:"room_name"`
}
type GameTableResolutionLogResource struct {
	Resolution int    `json:"resolution"`
	RoomName   string `json:"room_name"`
}
