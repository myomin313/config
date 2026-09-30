package requests

// PlatformApiKeyCreateRequest represents the request body for creating a Platform API Key
type PlatformApiKeyCreateRequest struct {
	PlatformID uint64 `json:"platform_id" validate:"required" example:"1"`                        // ID of the platform
	Status     string `json:"status" validate:"omitempty,oneof=active disabled" example:"active"` // Status of the API key; optional, defaults to "disabled"
}

// PlatformApiKeyUpdateRequest represents the request body for updating a Platform API Key
type PlatformApiKeyUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active disabled" example:"disabled"` // New status; must be "active" or "disabled"
}
