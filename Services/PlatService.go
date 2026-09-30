package Services

import (
	"config-service/Repositories"
	"config-service/models"
	"crypto/rand"
	"encoding/hex"
)

type PlatformApiKeyService struct {
	Repo *Repositories.PlatformApiKeyRepository
}

func generateSecret(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *PlatformApiKeyService) Create(platformID uint64) (*models.PlatformApiKey, error) {
	apiKey := &models.PlatformApiKey{
		PlatformID: platformID,
		Key:        generateSecret(16),
		Secret:     generateSecret(32),
		Status:     "enabled",
	}

	err := s.Repo.Create(apiKey)
	return apiKey, err
}
