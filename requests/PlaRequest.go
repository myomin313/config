package requests

// PlatformCreateRequest represents the request body for creating a Platform
type PlatformCreateRequest struct {
	Name        string   `json:"name" validate:"required" example:"Kla-klouk"`              // Name of the platform
	Email       string   `json:"email" validate:"required,email" example:"tech@stripe.com"` // Email of the platform
	Status      string   `json:"status"  example:"active"`                                  // Status; optional, defaults to "disabled"
	IPAllowlist []string `json:"ip_allowlist" example:"[\"127.0.0.1\"]"`                    // Optional JSON array of allowed IPs
	ApiKey      string   `json:"api_key" validate:"required"`
	SecretKey   string   `json:"secret_key" validate:"required"`
}

// PlatformUpdateRequest represents the request body for updating a Platform
type PlatformUpdateRequest struct {
	Name        string   `json:"name" example:"Kla-klouk"` // Updated name
	Email       string   `json:"email" validate:"required,email" example:"tech@stripe.com"`
	Status      string   `json:"status"  example:"active"` // Status; must be "active" or "disabled"
	IPAllowlist []string `json:"ip_allowlist" example:"[\"192.168.1.1\"]"`
	ApiKey      string   `json:"api_key"`
	SecretKey   string   `json:"secret_key"` // Updated IP allowlist
}
type CreateMerchantCurrencyRequest struct {
	MerchantCode string `json:"merchant_code" validate:"required"`
	Currency     string `json:"currency" validate:"required"`
	Status       string `json:"status" validate:"required,notblank,oneof=enable disabled" example:"enable"`
}
type UpdateMerchantCurrencyRequest struct {
	MerchantCode string `json:"merchant_code" validate:"required"`
	Currency     string `json:"currency" validate:"required"`
	Status       string `json:"status" validate:"required,notblank,oneof=enable disabled" example:"enable"`
}
