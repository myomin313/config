package Resources

type GameTypeResource struct {
	ID           uint64 `json:"id"`
	GameTypeName string `json:"game_type_name"`

	//GameTypeTranslations []GameTypeTranslationResource
}
