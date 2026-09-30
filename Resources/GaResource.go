package Resources

type GameResource struct {
	ID             uint64 `json:"id"`
	Status         string `json:"status"`
	BettingTime    int    `json:"betting_time"`
	SettlementTime int    `json:"settlement_time"`
	NoAutoExits    int    `json:"no_auto_exits"`
	CreatedBy      string `json:"created_by"` // now just string
	UpdatedBy      string `json:"updated_by"` // now just string
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	EnName         string `json:"en_name"`
	ZhName         string `json:"zh_name"`
	KmName         string `json:"km_name"`
	GameDataID     uint   `json:"game_data_id"`
	//GameTypeID     uint64 `json:"game_type_id"`
	//	GameTypeIDs []uint64 `json:"game_type_ids"`
	//Translations   []GameTranslationResource `json:"translations"`
	//GameTypes      []GameTypeResource        `json:"game_types"`
	//GameTypeData []GameTypeDataResource `json:"game_type_data"`
	GameTypeData []GameTypeGroupResource `json:"game_type_data"`
}

type GameTranslationResource struct {
	ID       uint64 `json:"id"`
	GameID   uint64 `json:"game_id"`
	Language string `json:"language_code"`
	Content  string `json:"content"`
}

// type GameTypeResource struct {
// 	ID   uint64 `json:"id"`
// 	GameID uint64 `json:game_id`
// 	Odds int `json:"odds"`
// 	Image string `json:"image"`

// 	EnName string `json:"en_name"`
// 	ZhName string `json:"zh_name"`
// 	KmName string `json:"km_name"`

//		//GameTypeTranslations []GameTypeTranslationResource
//	}
type GameTypeTranslationResource struct {
	ID               uint64 `json:"id"`
	Language         string `json:"language_code"`
	MasterGameTypeID uint64 `json:"master_game_type_id"`
	Name             string `json:"name"`
}

type GameTypeGroupResource struct {
	GameTypeID   uint64                 `json:"game_type_id"`
	GameTypeName string                 `json:"game_type_name"` // optional, can get from MasterGameType
	GameTypeData []GameTypeDataResource `json:"game_type_data"`
}

type GameTypeDataResource struct {
	ID               uint64  `json:"id"`
	Odds             float64 `json:"odds"`
	NameKm           string  `json:"name_km"`
	MasterGameTypeID uint64  `json:"master_game_type_id"`

	// 🔽 coming from MasterGameType
	Image    string `json:"image"`
	ImageUrl string `json:"image_url"`
	KeyName  string `json:"key_name"`
	NameEn   string `json:"name_en"`
	NameZh   string `json:"name_zh"`
}

type GamePlayTypeResource struct {
	ID           uint64 `json:"id"`
	GameTypeName string `json:"game_type_name"`
}

//	type GameTypeDataResource struct {
//		ID   uint64 `json:"id"`
//		Odds int `json:"odds"`
//		NameKm string `json:"name_km"`
//		MasterGameTypeID uint64 `json:"master_game_type_id"`
//		MasterGameType MasterGameTypeResource `json:"master_game_type"`
//	}

type GroupedMasterGameType struct {
	GameTypeID   uint64                   `json:"game_type_id"`
	GameTypeName string                   `json:"game_type_name"`
	Items        []MasterGameTypeResource `json:"items"`
}
type MasterGameTypeResource struct {
	ID      uint64 `json:"id"`
	Odds    int    `json:"odds"`
	Image   string `json:"image"`
	KeyName string `json:"key_name"`
	NameEn  string `json:"name_en"`
	NameKm  string `json:"name_km"`
	NameZh  string `json:"name_zh"`
}

type GameListResource struct {
	ID             uint64 `json:"id"`
	BettingTime    int    `json:"betting_time"`
	SettlementTime int    `json:"settlement_time"`
	//Name         string `json:"name"`
	EnName string `json:"en_name"`
	ZhName string `json:"zh_name"`
	KmName string `json:"km_name"`
}
type GameListFEResource struct {
	ID             uint64 `json:"id"`
	BettingTime    int    `json:"betting_time"`
	SettlementTime int    `json:"settlement_time"`
	NoAutoExits    int    `json:"no_auto_exits"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	//EnName string `json:"en_name"`
	//ZhName string `json:"zh_name"`
	//KmName string `json:"km_name"`
}
type LogGameResource struct {
	ID             uint64 `json:"id"`
	Status         string `json:"status"`
	BettingTime    int    `json:"betting_time"`
	SettlementTime int    `json:"settlement_time"`
	NoAutoExits    int    `json:"no_auto_exits"`
	CreatedBy      string `json:"created_by"` // now just string
	UpdatedBy      string `json:"updated_by"` // now just string
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	EnName         string `json:"en_name"`
	ZhName         string `json:"zh_name"`
	KmName         string `json:"km_name"`
	//GameTypeID     uint64 `json:"game_type_id"`
	//	GameTypeIDs []uint64 `json:"game_type_ids"`
	//Translations   []GameTranslationResource `json:"translations"`
	//GameTypes      []GameTypeResource        `json:"game_types"`
	//GameTypeData []GameTypeDataResource `json:"game_type_data"`
	GameTypeData []GameTypeLogGroupResource `json:"game_type_data"`
}
type GameTypeLogGroupResource struct {
	GameTypeID   uint64                    `json:"game_type_id"`
	GameTypeName string                    `json:"game_type_name"` // optional, can get from MasterGameType
	GameTypeData []GameTypeDataLogResource `json:"game_type_data"`
}
type GameTypeDataLogResource struct {
	//ID               uint64 `json:"id"`
	Odds             float64 `json:"odds"`
	NameKm           string  `json:"name_km"`
	MasterGameTypeID uint64  `json:"master_game_type_id"`

	// 🔽 coming from MasterGameType
	Image    string `json:"image"`
	ImageUrl string `json:"image_url"`
	KeyName  string `json:"key_name"`
	NameEn   string `json:"name_en"`
	NameZh   string `json:"name_zh"`
}
