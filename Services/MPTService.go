package Services

import (
	"config-service/Repositories"
	"config-service/models"
	// "config-service/requests"
)

type MasterPlayTypeService struct {
	Repo *Repositories.MasterPlayTypeRepository
}

// func (s *MasterGameTypeService) All(gameTypeID uint64) ([]models.MasterGameType, error) {
// 	return s.Repo.GetAll(gameTypeID)
// }

func (s *MasterPlayTypeService) All(playTypeIDs ...uint64) ([]models.MasterGameType, error) {
	return s.Repo.GetAll(playTypeIDs...)
}
