package requests

type GameHallCreateRequest struct {
	Status      string            `json:"status"`
	CreatedByID uint              `json:"created_by_id"`
	UpdatedByID uint              `json:"updated_by_id"`
	SortOrder   uint              `json:"sort_order" validate:"required"`
	Names       map[string]string `json:"names" validate:"required,require_all_game_locales"` // en, zh-CN, km
}
type GameHallUpdateRequest struct {
	Status      string            `json:"status"`
	UpdatedByID uint              `json:"updated_by_id"`
	SortOrder   uint              `json:"sort_order" validate:"required"`
	Names       map[string]string `json:"names" validate:"required,require_all_game_locales"` // en, zh-CN, km
}
