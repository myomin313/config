package requests

type GameRedLimitItem struct {
	//MinLimit uint `json:"min_limit" validate:"required,gte=0" example:"1"`
	//MaxLimit uint `json:"max_limit" validate:"required,gtefield=MinLimit" example:"100"`
	MinLimit float64 `json:"min_limit" validate:"required,gte=0" example:"1.50"`
	MaxLimit float64 `json:"max_limit" validate:"required,gtefield=MinLimit" example:"100.25"`
}

// RedLimitUpdateRequest represents the request body for creating a Red limit table
type RedLimitCreateRequest struct {
	GameTableName        map[string]string           `json:"game_table_name" validate:"required,require_all_game_locales" example:"{\"en\":\"Dragon Tiger\",\"km\":\"ដ្រាកន៍ ថាយហ្គ័\"}"`
	GameType             uint                        `json:"game" validate:"required,notblank" example:"2"`
	GameTypeCurrency     uint                        `json:"game_type_currency" validate:"required,notblank" example:"2"`
	GamePlayType         uint                        `json:"game_play_type" validate:"required,notblank" example:"2"`
	GamePlayTypeCurrency uint                        `json:"game_play_type_currency" validate:"required,notblank" example:"2"`
	Status               string                      `json:"status" validate:"required,notblank,oneof=active  disabled  maintenance" example:"active"`
	Sort                 uint                        `json:"sort" validate:"gt=0" example:"0"`
	GameRedLimit         map[string]GameRedLimitItem `json:"game_red_limit" validate:"required"`
}

// RedLimitUpdateRequest represents the request body for updating a Red limit table
type RedLimitUpdateRequest struct {
	GameTableName        *map[string]string           `json:"game_table_name" validate:"omitempty,required,require_all_game_locales" example:"{\"en\":\"Dragon Tiger\",\"km\":\"ដ្រាកន៍ ថាយហ្គ័\"}"`
	GameType             *uint                        `json:"game" validate:"omitempty,required,notblank" example:"2"`
	GameTypeCurrency     *uint                        `json:"game_type_currency" validate:"omitempty,required,notblank" example:"2"`
	GamePlayType         *uint                        `json:"game_play_type" validate:"omitempty,required,notblank" example:"2"`
	GamePlayTypeCurrency *uint                        `json:"game_play_type_currency" validate:"omitempty,required,notblank" example:"2"`
	Status               *string                      `json:"status" validate:"omitempty,required,notblank,oneof=active disabled maintenance" example:"disabled"`
	Sort                 *uint                        `json:"sort" validate:"omitempty,gt=0" example:"0"`
	GameRedLimit         *map[string]GameRedLimitItem `json:"game_red_limit" validate:"omitempty"`
}

// RedLimitUpdateRequest represents the request body for serching a Red limit table
type RedLimitSearchRequest struct {
	TableName  *string `query:"name" validate:"omitempty,notblank" example:"dragon-tiger"`
	CurrencyID *uint   `query:"currency_id" validate:"omitempty,gt=0" example:"1"`
}
