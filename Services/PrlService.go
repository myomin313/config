package Services

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"config-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type PlayRedLimitService struct {
	Repo     *Repositories.PlayRedLimitRepository
	PlayRepo *Repositories.PlayRepository
	Audit    *helper.AuditService
	Redis    *redis.Client
}

// Create a new Platform
func (s *PlayRedLimitService) CreateRedLimit(req *requests.RedLimitCreateRequest, id uint) (*models.GameRedLimit, error) {
	playtable := models.GameRedLimit{
		GameType:             req.GameType,
		GameTypeCurrency:     req.GameTypeCurrency,
		GamePlayType:         req.GamePlayType,
		GamePlayTypeCurrency: req.GamePlayTypeCurrency,
		Status:               req.Status,
		Sort:                 req.Sort,
		CreatedByID:          id,
		UpdatedByID:          id,
	}

	// Red limit translations
	for lang, name := range req.GameTableName {
		playtable.Translations = append(playtable.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "red_limit",
			Content:      name,
			Type:         "api",
			TableType:    "game_red_limits",
		})
	}

	// Game types
	for key, gtd := range req.GameRedLimit {
		playRedlimit := models.RedLimitTable{
			MinLimit:        gtd.MinLimit,
			MaxLimit:        gtd.MaxLimit,
			GameRedLimitKey: key,
		}
		playtable.RedLimitTable = append(playtable.RedLimitTable, playRedlimit)
	}

	if err := s.Repo.Create(&playtable); err != nil {
		return nil, err
	}

	return &playtable, nil
}

// All returns all Platforms
func (s *PlayRedLimitService) AllRedLimits(req *requests.RedLimitSearchRequest, lang string) ([]models.GameRedLimit, error) {
	var platforms []models.GameRedLimit

	if err := s.Repo.FindAll(req, &platforms, lang); err != nil {
		return nil, err
	}
	for i := range platforms {
		sortRedLimitTable(platforms[i].RedLimitTable)
	}
	return platforms, nil
}

// FindByID returns a Platform by ID
func (s *PlayRedLimitService) FindRedLimitByID(id uint64) (*models.GameRedLimit, error) {
	var platform models.GameRedLimit
	err := s.Repo.FindByID(id, &platform)
	if err != nil {
		return nil, errors.New("record not found")
	}
	sortRedLimitTable(platform.RedLimitTable)
	return &platform, nil
}

// Update a Platform by ID
func (s *PlayRedLimitService) UpdateRedLimit(id uint64, req *requests.RedLimitUpdateRequest, userId uint, lang string) (*models.GameRedLimit, error) {
	platform, err := s.FindRedLimitByID(id)
	if err != nil {
		return nil, err
	}

	// 🔥 Apply only non-nil fields from request to model
	utils.UpdateStruct(platform, req)

	platform.UpdatedByID = userId
	platform.UpdatedAt = time.Now()

	if req.GameTableName != nil && len(*req.GameTableName) > 0 {
		// ---- Replace translations ----
		if errR := s.PlayRepo.DeletePlayTranslations(id, "game_red_limits"); errR != nil {
			return nil, errR
		}
		platform.Translations = nil
		for lang, name := range *req.GameTableName {
			platform.Translations = append(platform.Translations, models.GameTranslation{
				LanguageCode: lang,
				LangName:     "red_limit",
				Content:      name,
				Type:         "api",
				TableType:    "game_red_limits",
			})
		}

		//log.Println("platform.Translations", platform.Translations)
	}

	//update red limit table
	if req.GameRedLimit != nil && len(*req.GameRedLimit) > 0 {
		if errRed := s.Repo.DeleteRedlimitTable(id); errRed != nil {
			return nil, errRed
		}
		platform.RedLimitTable = nil
		for key, gtd := range *req.GameRedLimit {
			playRedlimit := models.RedLimitTable{
				MinLimit:        gtd.MinLimit,
				MaxLimit:        gtd.MaxLimit,
				GameRedLimitKey: key,
			}
			platform.RedLimitTable = append(platform.RedLimitTable, playRedlimit)
		}
	}
	if errP := s.Repo.Update(platform); errP != nil {
		log.Println("game table update error", errP)
		return nil, errP
	}

	return platform, nil
}

// Delete a Platform by ID
func (s *PlayRedLimitService) DeleteRedLimit(id uint64) error {
	log.Println("red limit DeleteRedLimit")
	platform, err := s.FindRedLimitByID(id)
	if err != nil {
		return err
	}
	if err := s.PlayRepo.DeletePlayTranslations(id, "game_red_limits"); err != nil {
		return fmt.Errorf("failed to delete game table translations: %w", err)
	}

	if err := s.Repo.DeleteRedlimitTable(id); err != nil {
		return fmt.Errorf("failed to delete red limit table: %w", err)
	}

	if err := s.Repo.Delete(platform); err != nil {
		return fmt.Errorf("failed to delete game table: %w", err)
	}
	return nil
}
func (s *PlayRedLimitService) AuditAction(
	c *fiber.Ctx,
	action string,
	model string,
	user jwt.MapClaims,
	oldData any,
	newData any,
) {

	helper.PublishAudit(
		s.Audit,
		c,
		action,
		model,
		user,
		oldData,
		newData,
	)
}
func sortRedLimitTable(rows []models.RedLimitTable) {
	sort.SliceStable(rows, func(i, j int) bool {
		return utils.PlayOrderIndex(rows[i].GameRedLimitKey) <
			utils.PlayOrderIndex(rows[j].GameRedLimitKey)
	})
}
