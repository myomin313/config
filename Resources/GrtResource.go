package Resources

type RedLimit struct {
	MinLimit float64 `json:"min_limit"`
	MaxLimit float64 `json:"max_limit"`
}
type GameRedLimitResource struct {
	ID                   uint                `json:"id"`
	EnName               string              `json:"en_name"`
	ChName               string              `json:"ch_name"`
	KhName               string              `json:"kh_name"`
	GameType             uint                `json:"game"`
	GameTypeCurrency     uint                `json:"game_type_currency"`
	GamePlayType         uint                `json:"game_play_type"`
	GamePlayTypeCurrency uint                `json:"game_play_type_currency"`
	Status               string              `json:"status"`
	CreatedBy            string              `json:"created_by"`
	Sort                 uint                `json:"sort"`
	UpdatedBy            string              `json:"updated_by"`
	CreatedAt            string              `json:"created_at"`
	UpdatedAt            string              `json:"updated_at"`
	GameRedLimit         map[string]RedLimit `json:"game_red_limit,omitempty"` // omit if nil
}
type GameRedLimitLogResource struct {
	ID               uint                `json:"id"`
	EnName           string              `json:"en_name"`
	ChName           string              `json:"ch_name"`
	KhName           string              `json:"kh_name"`
	GameTypeCurrency string              `json:"game_type_currency"`
	GamePlayType     string              `json:"game_play_type"`
	Status           string              `json:"status"`
	CreatedBy        string              `json:"created_by"`
	Sort             uint                `json:"sort"`
	UpdatedBy        string              `json:"updated_by"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
	GameRedLimit     map[string]RedLimit `json:"game_red_limit,omitempty"` // omit if nil
	Game             *GameLogResource    `json:"game,omitempty"`
}
