package Resources

type AnnouncementResource struct {
	ID               uint   `json:"id"`
	AnnouncementName string `json:"announcement_name"`
	Status           string `json:"status"`
	//Title          string `json:"title"`
	SortOrder          uint                         `json:"sort_order"`
	BroadcastAudience  string                       `json:"broadcast_audience"`
	CreatedBy          string                       `json:"created_by"` // now just string
	UpdatedBy          string                       `json:"updated_by"` // now just string
	CreatedAt          string                       `json:"created_at"`
	UpdatedAt          string                       `json:"updated_at"`
	AnnouncementDetail []AnnouncementDetailResource `json:"announcement_detail"`
}
type AnnouncementDetailResource struct {
	ID           uint   `json:"id"`
	LanguageCode string `json:"language_code"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}
type AnnouncementListFEResource struct {
	ID                uint   `json:"id"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Status            string `json:"status"`
	SortOrder         uint   `json:"sort_order"`
	BroadcastAudience string `json:"broadcast_audience"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type AnnouncementLogResource struct {
	ID                uint   `json:"id"`
	AnnouncementName  string `json:"announcement_name"`
	Status            string `json:"status"`
	SortOrder         uint   `json:"sort_order"`
	BroadcastAudience string `json:"broadcast_audience"`
	CreatedBy         string `json:"created_by"` // now just string
	UpdatedBy         string `json:"updated_by"` // now just string
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
type AnnouncementDetailLogResource struct {
	ID               uint   `json:"id"`
	LanguageCode     string `json:"language_code"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	AnnouncementName string `json:"announcement_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
