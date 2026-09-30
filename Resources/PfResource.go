package Resources

type PlatformResource struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// MerchantCode string   `json:"merchant_code"`
	Email       string   `json:"email"`
	Status      string   `json:"status"`
	IPAllowlist []string `json:"ip_allowlist"`
	ApiKey      string   `json:"api_key"`
	SecretKey   string   `json:"secret_key"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
type PlatformsResource struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CurrencyResource struct {
	ID       uint   `json:"id"`
	Currecny string `json:"currency"`
}

type MerchantResource struct {
	ID           uint   `json:"id"`
	Currecny     string `json:"currency"`
	MerchantCode string `json:"merchant_code"`
	Name         string `json:"name"`
	PlatformID   uint   `json:"platform_id"`
	Email        string `json:"email"`
	Status       string `json:"status"`
}
type MerchantCurrencyLogResource struct {
	ID           uint   `json:"id"`
	Currency     string `json:"currency"`
	MerchantCode string `json:"merchant_code"`
	Status       string `json:"status"`
}
