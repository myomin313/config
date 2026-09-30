package Services

import (
	"config-service/Repositories"
	"config-service/models"
)

type PlayDataService struct {
	Repo *Repositories.PlayDataRepository
}

func (s *PlayDataService) All() ([]models.GameData, error) {
	return s.Repo.FindAll()
}
