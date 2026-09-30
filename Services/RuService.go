package Services

import (
	"fmt"
	"os"
)

type RuleService struct {
}

// GetRuleImageURL constructs the MinIO URL for rule images based on game parameters
func (rs *RuleService) GetRuleImageURL(playID, playType, agent, languageCode string) (string, error) {
	// Get MinIO endpoint from environment
	minioEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if minioEndpoint == "" {
		return "", fmt.Errorf("MINIO_PUBLIC_ENDPOINT not configured")
	}
	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		return "", fmt.Errorf("MINIO_BUCKET not configured")
	}

	// Construct the rule image URL
	// Format: MINIO_PUBLIC_ENDPOINT + '/' + 'kh-game' + '/rule_' + game_id + '_' + game_type + '_' + agent + '_' + language_code + '.png'
	url := fmt.Sprintf("%s/%s/kh-game/rule_%s_%s_%s_%s.png", minioEndpoint, minioBucket, playID, playType, agent, languageCode)

	return url, nil
}
