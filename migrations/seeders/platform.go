package seeders

import (
	"config-service/config"
	"config-service/models"
	"crypto/rand"
	"encoding/hex"
	"log"

	"gorm.io/datatypes"
)

// generateSecret creates a random hex string of given byte length
func generateSecret(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random secret: " + err.Error())
	}
	return hex.EncodeToString(b)
}

func CreatePlatforms() {
	log.Println("=== Starting Platform Migration ===")

	platforms := []models.Platform{
		{
			ID:           1,
			MerchandCode: "PLT001",
			Name:         "Platform One",
			Email:        "klaklouk@example.com",
			Status:       "enable",
			IPAllowlist:  datatypes.JSONSlice[string]{}, // empty = allow all IPs
		},
	}

	for _, p := range platforms {
		var existing models.Platform

		log.Printf("Processing Platform ID: %d | Name: %s\n", p.ID, p.Name)

		// Check if platform already exists
		if err := config.DB.First(&existing, p.ID).Error; err == nil {
			log.Printf("Platform ID %d already exists, skipping...\n", p.ID)
			continue
		}

		// Create platform
		if err := config.DB.Create(&p).Error; err != nil {
			log.Printf("Failed to create Platform ID %d: %v\n", p.ID, err)
			continue
		}
		log.Printf("Platform Created | ID: %d | Name: %s\n", p.ID, p.Name)

		// Create API Key for platform with random key + secret
		apiKey := models.PlatformApiKey{
			PlatformID: p.ID,
			Status:     "active",
			Key:        "4fdb4b098b6c1a3ce5c6f4b88acde855",                                 // 32 hex chars
			Secret:     "343543ef9ba3b6e1f7d505f3eaa9fe9a395d9333dbb535926ce7bc2ced4c5cd8", // 64 hex chars
		}

		if err := config.DB.Create(&apiKey).Error; err != nil {
			log.Printf("Failed to create API Key for Platform ID %d: %v\n", p.ID, err)
		} else {
			log.Printf("API Key Created | PlatformID: %d | Key: %s\n", p.ID, apiKey.Key)
		}
	}

	log.Println("=== Platform Migration Completed ===")
}
