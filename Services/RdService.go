package Services

import (
	"config-service/Repositories"

	"gorm.io/gorm"
)

type RoundService struct {
	Repo *Repositories.RoundRepository
}

// NewRoundService creates a new RoundService
func NewRoundService(db *gorm.DB) *RoundService {
	return &RoundService{
		Repo: &Repositories.RoundRepository{
			DB: db,
		},
	}
}

// GetGameResults returns game result report data
func (s *RoundService) GetPlayResults() ([]Repositories.PlayResultRow, error) {
	var results []Repositories.PlayResultRow
	err := s.Repo.GetPlayResults(&results)
	return results, err
}
