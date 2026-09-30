package requests

import "config-service/helper"

type BannerImage struct {
	ID           *uint64 `json:"id"`
	ImagePath    *string `json:"image_path" validate:"omitempty,required"`
	LangCode     string  `json:"lang_code" validate:"required"`
	DefaultImage *bool   `json:"default_image" validate:"required"`
}
type BannerCreateRequest struct {
	ImageName       string             `json:"image_name" validate:"required,min=3"`
	DisplayPlatform string             `json:"display_platform" validate:"required,oneof=h5 web mobile"`
	DisplayDuration string             `json:"display_duration" validate:"required,oneof=schedule permanent"`
	StartDate       *helper.CustomTime `json:"start_date"`
	EndDate         *helper.CustomTime `json:"end_date"`
	BannerImage     []BannerImage      `json:"banner_images" validate:"required,dive"`
}

type BannerUpdateRequest struct {
	ImageName       *string            `json:"image_name" validate:"required,min=3"`
	DisplayPlatform *string            `json:"display_platform" validate:"omitempty,oneof=h5 web mobile"`
	DisplayDuration *string            `json:"display_duration" validate:"omitempty,oneof=schedule permanent"`
	StartDate       *helper.CustomTime `json:"start_date"`
	EndDate         *helper.CustomTime `json:"end_date"`
	BannerImage     []BannerImage      `json:"banner_images" validate:"omitempty,dive"`
}

type BannerSearchRequest struct {
	ImageName       *string `query:"image_name" validate:"omitempty,notblank" example:"dragon-tiger"`
	DisplayPlatform *string `query:"display_platform" validate:"omitempty,required,oneof=h5 web mobile"`
	DisplayDuration *string `query:"display_duration" validate:"omitempty,required,oneof=schedule permanent"`
	Status          *string `query:"status" validate:"omitempty,oneof=active  disabled"`
}

type BannerMassStatusUpdateRequest struct {
	BannerIDs []uint `json:"banner_ids" validate:"required,min=1" example:"1,2,3"`
	Status    string `json:"status" validate:"required,oneof=active  disabled" example:"active"`
}

type BannerStatusUpdateRequest struct {
	BannerID uint   `json:"banner_id" validate:"required,min=1" example:"1"`
	Status   string `json:"status" validate:"required,oneof=active  disabled" example:"active"`
}

type BannerMassDeleteRequest struct {
	BannerIDs []uint `json:"banner_ids" validate:"required,min=1" example:"1,2,3"`
}

type GetBannerRequest struct {
	DisplayPlatform string `json:"display_platform" validate:"required,oneof=h5 web mobile"`
}
