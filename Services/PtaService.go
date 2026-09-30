package Services

import (
	"config-service/models"
	"config-service/pb/commonpb"

	"gorm.io/gorm"
)

type PlayTableActiveService struct {
	commonpb.UnimplementedGameTableActiveServiceServer
	DB *gorm.DB
}

// GetTables implements the streaming RPC you defined
func (s *PlayTableActiveService) GetActiveTables(req *commonpb.GameTableRequest, stream commonpb.GameTableActiveService_GetActiveTablesServer) error {
	var tables []models.GameTable

	// Start query
	query := s.DB.Model(&models.GameTable{})

	// Always filter ACTIVE
	query = query.Where("status = ?", "disable")
	if req.GameType != 0 {
		query = query.Where("game_type = ?", req.GameType)
	}

	// Optional filter: AffiliateOffice
	if req.AffiliateOffice != nil {
		query = query.Where("affiliate_office = ?", *req.AffiliateOffice)
	}

	// Optional filter: GameType
	//if req.GameType{

	//}

	// Execute query
	if err := query.Find(&tables).Error; err != nil {
		return err
	}

	for _, t := range tables {
		// Map your DB model to your Proto message
		resp := &commonpb.GetActiveTablesResponse{
			Id:              uint64(t.ID),
			Mode:            t.Mode, //"digital", // Example field, map as needed
			BettingTime:     uint32(t.BettingTime),
			SettlementTime:  uint32(t.SettlementTime),
			Status:          t.Status, // Example field, map as needed
			AffiliateOffice: uint32(t.AffiliateOffice),
			GameType:        uint32(t.GameType),
			TableNumber:     t.TableNumber,
			MobileIcon:      t.MobileIcon,
			//TableResultType: t.TableResultType,
			Sort:   uint32(t.Sort),
			PcIcon: t.PCIcon,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}
