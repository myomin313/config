package Services

import (
	"config-service/models"
	"config-service/pb/commonpb"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type GrpcPlayService struct {
	DB *gorm.DB
	commonpb.UnimplementedGameServiceServer
}

func (s *GrpcPlayService) GetGamesByIds(
	ctx context.Context,
	req *commonpb.GameIdsRequest,
) (*commonpb.GameListResponse, error) {

	var translations []models.GameTranslation

	err := s.DB.
		Where("table_id IN ?", req.Ids).
		Where("language_code = ?", req.LanguageCode).
		Where("lang_name = ?", "game_name").
		Find(&translations).Error

	fmt.Println("REQ IDS:", req.Ids)
	fmt.Println("LANG:", req.LanguageCode)
	fmt.Println("TRANSLATIONS:", translations)

	if err != nil {
		return nil, err
	}

	plays := make([]*commonpb.GameItem, 0)

	for _, t := range translations {
		plays = append(plays, &commonpb.GameItem{
			Id:   t.TableID,
			Name: t.Content,
		})
	}

	return &commonpb.GameListResponse{
		Games: plays,
	}, nil
}

func (s *GrpcPlayService) GetHallNames(
	ctx context.Context,
	req *commonpb.HallNamesRequest,
) (*commonpb.HallListResponse, error) {

	var translations []models.GameTranslation

	err := s.DB.
		Where("table_id IN ?", req.Ids).
		Where("language_code = ?", req.LanguageCode).
		Where("lang_name = ?", "game_hall_name").
		Find(&translations).Error

	fmt.Println("Hall REQ IDS:", req.Ids)
	fmt.Println("Hall LANG:", req.LanguageCode)
	fmt.Println("Hall TRANSLATIONS:", translations)

	if err != nil {
		return nil, err
	}

	halls := make([]*commonpb.HallItem, 0)

	for _, t := range translations {
		halls = append(halls, &commonpb.HallItem{
			Id:   t.TableID,
			Name: t.Content,
		})
	}

	return &commonpb.HallListResponse{
		Halls: halls,
	}, nil
}
