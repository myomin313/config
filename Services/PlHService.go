package Services

import (
	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"fmt"

	"config-service/requests"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type PlayHallService struct {
	Repo     *Repositories.PlayHallRepository
	PlayRepo *Repositories.PlayRepository
	Audit    *helper.AuditService
	Redis    *redis.Client
}

func (s *PlayHallService) CreatePlayHall(req requests.GameHallCreateRequest, id uint) (*models.GameHall, error) {

	playhall := &models.GameHall{
		Status:      req.Status,
		CreatedByID: id,
		UpdatedByID: id,
		SortOrder:   req.SortOrder,
	}

	// for lang, name := range req.Names {
	// 	gamehall.Translations = append(gamehall.Translations, models.GameTranslation{
	// 		LanguageCode: lang,
	// 		LangName:     "game_hall_name",
	// 		Content:      name,
	// 		TableType:    "game_halls",
	// 		Type:         "api",
	// 	})
	// }

	for lang, name := range req.Names {
		if utf8.RuneCountInString(name) > 15 {
			return nil, fmt.Errorf("name for language %s exceeds 15 characters", lang)
		}

		playhall.Translations = append(playhall.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "game_hall_name",
			Content:      name,
			TableType:    "game_halls",
			Type:         "api",
		})
	}

	if err := s.Repo.Create(playhall); err != nil {
		return nil, err
	}

	return s.Repo.FindByID(playhall.ID)
}
func (s *PlayHallService) All() ([]models.GameHall, error) {
	return s.Repo.FindAll()
}
func (s *PlayHallService) AllEnable() ([]models.GameHall, error) {
	return s.Repo.FindAllEnable()
}
func (s *PlayHallService) Delete(id uint64) error {
	playhall, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.PlayRepo.DeletePlayTranslations(id, "game_halls"); err != nil {
		return fmt.Errorf("failed to delete game hall translations: %w", err)
	}

	if err := s.Repo.Delete(playhall); err != nil {
		return fmt.Errorf("failed to delete game hall: %w", err)
	}

	tables, err := helper.GetTablesByIDs([]uint{}, uint(id), 0)
	if err == nil {
		BulkUpdateTableInfo(tables, "disable")
	}

	return nil

}
func (s *PlayHallService) Update(playhallID uint64, req requests.GameHallUpdateRequest, userId uint) (*models.GameHall, error) {

	// Load existing game
	playhall, err := s.Repo.FindByID(playhallID)
	if err != nil {
		return nil, err
	}

	playhall.Status = req.Status
	playhall.UpdatedByID = userId
	playhall.SortOrder = req.SortOrder

	if err := s.PlayRepo.DeletePlayTranslations(playhallID, "game_halls"); err != nil {
		return nil, err
	}
	playhall.Translations = nil
	for lang, name := range req.Names {
		if utf8.RuneCountInString(name) > 15 {
			return nil, fmt.Errorf("name for language %s exceeds 15 characters", lang)
		}
		playhall.Translations = append(playhall.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "game_hall_name",
			Content:      name,
			TableType:    "game_halls",
			Type:         "api",
		})
	}

	// Save game + associations
	if err := s.Repo.Update(playhall); err != nil {
		return nil, err
	}
	tables, err := helper.GetTablesByIDs([]uint{}, uint(playhallID), 0)
	if err == nil {
		BulkUpdateTableInfo(tables, playhall.Status)
	}
	return s.Repo.FindByID(playhallID)
	//return game, nil
}
func (s *PlayHallService) AuditAction(
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
	// if s.Audit == nil {
	// 	log.Println("❌ Audit service not initialized")
	// 	return
	// }
	// if s.Audit.Redis == nil {
	// 	log.Println("❌ Redis client not initialized")
	// 	return
	// }
	// ua := user_agent.New(c.Get("User-Agent"))
	// device := "Desktop"
	// if ua.Mobile() {
	// 	device = "Mobile"
	// }
	// log.Println("role:", user["role"])
	// // ✅ DTO (NOT model)
	// audit := map[string]any{
	// 	"action":     action,
	// 	"model":      model,
	// 	"ip_address": c.IP(),
	// 	"device":     device,
	// 	"created_at": time.Now(),
	// 	"role":       user["role"],
	// 	"user": map[string]any{
	// 		"email": user["email"],
	// 	},
	// }
	// if before != nil {
	// 	audit["old_data"] = before
	// }
	// if after != nil {
	// 	audit["new_data"] = after
	// }
	// s.Audit.Publish(model, audit, time.Minute)
}
func (s *PlayHallService) FindByID(id uint64) (*models.GameHall, error) {
	return s.Repo.FindByID(id)
}
