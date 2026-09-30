package Services

import (
	"context"

	"config-service/models"
	"config-service/pb/commonpb"

	"gorm.io/gorm"
)

type GrpcUserService struct {
	commonpb.UnimplementedPlatformServiceServer
	DB *gorm.DB
}

// GetPlatformByCode (UNARY RPC)
func (s *GrpcUserService) GetPlatform(ctx context.Context, req *commonpb.PlatformEmpty) (*commonpb.PlatformResponse, error) {

	var platform models.Platform

	err := s.DB.
		Preload("APIKey").
		Preload("MerchantCurrencies").
		//Where("merchand_code = ?", req.MerchantCode).
		First(&platform).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return &commonpb.PlatformResponse{
			Found: false,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	// Map DB → Proto
	resp := &commonpb.PlatformResponse{
		Id:           platform.ID,
		MerchantCode: platform.MerchandCode,
		Name:         platform.Name,
		Email:        platform.Email,
		Status:       platform.Status,
		IpAllowlist:  platform.IPAllowlist,
		Found:        true,
	}

	// Attach API Key if exists
	if platform.APIKey.ID != 0 {
		resp.ApiKey = &commonpb.PlatformApiKeyMessage{
			Id:     platform.APIKey.ID,
			Key:    platform.APIKey.Key,
			Secret: platform.APIKey.Secret,
			Status: platform.APIKey.Status,
		}
	}
	for _, mc := range platform.MerchantCurrencies {
		resp.MerchantCurrencies = append(
			resp.MerchantCurrencies,
			&commonpb.MerchantCurrencyMessage{
				Id:           uint64(mc.ID),
				MerchantCode: mc.MerchantCode,
				Currency:     mc.Currency,
			},
		)
	}

	return resp, nil
}
