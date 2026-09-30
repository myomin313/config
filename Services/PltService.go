package Services

import (
	"config-service/Repositories"
	"config-service/models"
	// "config-service/requests"
)

type PlayTypeService struct {
	Repo *Repositories.PlayTypeRepository
}

func (s *PlayTypeService) All(play_data_id uint) ([]models.GameType, error) {
	return s.Repo.GetAll(play_data_id)
}
func (s *PlayTypeService) FindByID(play_type_id uint64) (*models.GameType, error) {
	return s.Repo.FindByID(play_type_id)
}
