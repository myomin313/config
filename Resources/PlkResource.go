package Resources

type PlatformApiKeyResource struct {
	ID         uint64 `json:"id"`
	PlatformID uint64 `json:"platform_id"`
	Status     string `json:"status"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
	CreatedAt  string `json:"created_at"`
}
