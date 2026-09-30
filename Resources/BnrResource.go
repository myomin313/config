package Resources

type BannerImageResource struct {
	ID           uint64 `json:"id"`
	ImagePath    string `json:"image_path"`
	LangCode     string `json:"lang_code"`
	DefaultImage *bool  `json:"default_image"`
}
type BannerResource struct {
	ID              uint64                `json:"id"`
	ImageName       string                `json:"image_name"`
	DisplayPlatform string                `json:"display_platform"`
	DisplayDuration string                `json:"display_duration"`
	Status          string                `json:"status"`
	StartDate       string                `json:"start_date"`
	EndDate         string                `json:"end_date"`
	CreatedBy       string                `json:"created_by"`
	UpdatedBy       string                `json:"updated_by"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
	BannerImages    []BannerImageResource `json:"banner_images,omitempty"`
}

type BannerListResource struct {
	ImageName    string                `json:"image_name"`
	BannerImages []BannerImageResource `json:"banner_images,omitempty"`
}
