package requests

type AnnouncementCreateRequest struct {
	AnnouncementName string `json:"announcement_name" validate:"required"`
	SortOrder    uint `json:"sort_order" validate:"required"`
	BroadcastAudience string `json:"broadcast_audience" validate:"required,oneof=member agent all" example:"member"`
	Status         string `json:"status" validate:"required"`
//	AnnouncementDetail []AnnouncementDetailRequest `json:"announcement_detail" validate:"required"`
}

type AnnouncementDetailRequest struct {
	//AnnouncementID         uint              `json:"announcement_id" validate:"required"`
	LanguageCode           string            `json:"language_code" validate:"required,oneof=en ch kh"`
	Title                  string            `json:"title" validate:"required"`
	Content                string            `json:"content" validate:"required"`
	//CreatedAt              string            `json:"created_at"`
	//UpdatedAt              string            `json:"updated_at"`
}
type AnnouncementUpdateRequest struct {
    AnnouncementName string `json:"announcement_name" validate:"required"`
	SortOrder    uint `json:"sort_order" validate:"required"`
	BroadcastAudience string `json:"broadcast_audience" validate:"required,oneof=member agent all" example:"member"`
	Status         string `json:"status" validate:"required"`
	//CreatedBy     uint `json:"created_by"`
   // UpdatedBy     uint `json:"updated_by"`
	//AnnouncementDetail []AnnouncementDetailRequest `json:"announcement_detail" validate:"required"`
}
type AnnouncementMassStatusUpdateRequest struct {
	AnnouncementIDs []uint `json:"announcement_ids" validate:"required,min=1" example:"1,2,3"`
	Status   string `json:"status" validate:"required" example:"enable"`
}
type DeleteAnnouncementsRequest struct {
	IDs []uint `json:"ids"`
}
type AnnouncementSearchRequest struct {
	AnnouncementName                  *string            `query:"announcement_name"`
	BroadcastAudience      *string            `query:"broadcast_audience"`
	Status                 *string            `query:"status"`
	LanguageCode           *string            `query:"language_code"`
}

