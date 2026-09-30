package Services

import (
	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	// "fmt"
)

type PlayService struct {
	Repo  *Repositories.PlayRepository
	Audit *helper.AuditService
	Redis *redis.Client
}

func (s *PlayService) All() ([]models.Game, error) {
	return s.Repo.FindAll()
}
func (s *PlayService) AllEnable() ([]models.Game, error) {
	return s.Repo.FindAllEnable()
}
func (s *PlayService) CreatePlay(req requests.GameCreateRequest, id uint) (*models.Game, error) {

	// exists, err := s.Repo.GameTypeExists(req.GameTypeID)
	// if err != nil {
	// 	return nil, err
	// }
	// if !exists {
	// 	return nil, errors.New("game_type_id does not exist")
	// }
	// return nil, errors.New("test error from service")

	existingPlay, err := s.Repo.FindByPlayDataID(req.GameDataID)

	if err == nil && existingPlay != nil {
		// 👉 Game already exists → UPDATE instead of CREATE

		updateReq := requests.GameUpdateRequest{
			GameTypeIDs:    req.GameTypeIDs,
			GameTypeData:   req.GameTypeData,
			Names:          req.Names,
			BettingTime:    req.BettingTime,
			SettlementTime: req.SettlementTime,
			NoAutoExits:    req.NoAutoExits,
			Status:         req.Status,
			GameDataID:     req.GameDataID,
		}

		return s.UpdatePlay(existingPlay.ID, updateReq, id)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var playTypes []models.GameType

	for _, id := range req.GameTypeIDs {

		exists, err := s.Repo.PlayTypeExists(id)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("game_type_id does not exist")
		}

		playTypes = append(playTypes, models.GameType{ID: id})
	}
	for _, gtd := range req.GameTypeData {
		exists, err := s.Repo.MasterPlayTypeExists(gtd.MasterGameTypeID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("master_game_type_id  does not exist")
		}
	}

	play := models.Game{
		BettingTime:    req.BettingTime,
		SettlementTime: req.SettlementTime,
		NoAutoExits:    req.NoAutoExits,
		Status:         req.Status,
		GameDataID:     req.GameDataID,
		ID:             uint64(req.GameDataID),
		GameTypes:      playTypes,
		CreatedBy:      id,
		UpdatedBy:      id,
	}

	// Game translations
	for lang, name := range req.Names {
		play.Translations = append(play.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "game_name",
			Content:      name,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	// Game types
	for _, gtd := range req.GameTypeData { // make sure your request struct has GameTypeDataRequest slice
		playTypeData := models.GameTypeData{
			MasterGameTypeID: gtd.MasterGameTypeID,
			Odds:             gtd.Odds,
			NameKm:           gtd.NameKm,
			Image:            gtd.Image,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// Link to parent game after we have ID
		// For GORM, we can use association later, but here we append to GameTypeData slice
		play.GameTypeData = append(play.GameTypeData, playTypeData)
	}

	if err := s.Repo.Create(&play); err != nil {
		return nil, err
	}
	return s.Repo.FindByID(play.ID)
	//return &game, nil
}

func (s *PlayService) UpdatePlay(playID uint64, req requests.GameUpdateRequest, userId uint) (*models.Game, error) {

	// Load existing game
	play, err := s.Repo.FindByID(playID)
	if err != nil {
		return nil, err
	}

	// exists, err := s.Repo.GameTypeExists(req.GameTypeID)
	// if err != nil {
	// 	return nil, err
	// }
	// if !exists {
	// 	return nil, errors.New("game_type_id does not exist")
	// }
	var playTypes []models.GameType

	for _, id := range req.GameTypeIDs {

		exists, err := s.Repo.PlayTypeExists(id)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("game_type_id does not exist")
		}

		playTypes = append(playTypes, models.GameType{ID: id})
	}

	for _, gtd := range req.GameTypeData {
		exists, err := s.Repo.MasterPlayTypeExists(gtd.MasterGameTypeID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("master_game_type_id  does not exist")
		}
	}
	// Update main fields
	play.BettingTime = req.BettingTime
	play.SettlementTime = req.SettlementTime
	play.NoAutoExits = req.NoAutoExits
	play.Status = req.Status
	play.GameDataID = req.GameDataID
	play.UpdatedBy = userId
	//game.GameTypeID = req.GameTypeID

	// This replaces records in game_game_types table
	if err := s.Repo.ReplacePlayTypes(play, playTypes); err != nil {
		return nil, err
	}
	// ---- Replace translations ----
	if err := s.Repo.DeletePlayTranslations(playID, "games"); err != nil {
		return nil, err
	}

	play.Translations = nil
	for lang, name := range req.Names {
		play.Translations = append(play.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "game_name",
			Content:      name,
			UpdatedAt:    time.Now(),
		})
	}

	// ---- Replace game types ----
	if err := s.Repo.DeletePlayTypeData(playID); err != nil {
		return nil, err
	}
	play.GameTypeData = []models.GameTypeData{}
	// Game types
	for _, gtd := range req.GameTypeData { // make sure your request struct has GameTypeDataRequest slice
		playTypeData := models.GameTypeData{
			MasterGameTypeID: gtd.MasterGameTypeID,
			Odds:             gtd.Odds,
			NameKm:           gtd.NameKm,
			Image:            gtd.Image,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// Link to parent game after we have ID
		// For GORM, we can use association later, but here we append to GameTypeData slice
		play.GameTypeData = append(play.GameTypeData, playTypeData)
	}
	// Save game + associations
	if err := s.Repo.Update(play); err != nil {
		return nil, err
	}
	tables, err := helper.GetTablesByIDs([]uint{}, 0, uint(playID))
	if err == nil {
		BulkUpdateTableInfo(tables, play.Status)
	}
	return s.Repo.FindByID(playID)
	//return game, nil
}

// func (s *GameService) DeleteGame(id uint64) error {
// 		_, err := s.Repo.FindByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	// 2️⃣ Delete game FIRST
// 	s.Repo.Delete(id);

// 	// 3️⃣ Delete game_type_data SECOND
// 	s.Repo.DeleteGameTypeData(id);
// 	// 4️⃣ Delete ONLY game_name translation LAST
// 	 s.Repo.DeleteGameTranslations(id,"game_name");

// 	return nil
// }

func (s *PlayService) DeletePlay(id uint64) error {
	result := s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		repo := &Repositories.PlayRepository{DB: tx}

		// 1️⃣ Find the game with GameTypes loaded
		var play models.Game
		if err := tx.Preload("GameTypes").First(&play, id).Error; err != nil {
			return err
		}

		// 2️⃣ Clear pivot table (game_game_types)
		if err := tx.Model(&play).Association("GameTypes").Clear(); err != nil {
			return err
		}
		// 4️⃣ Delete the game

		// 2️⃣ Delete the game_name translation
		if err := repo.DeletePlayTranslations(id, "games"); err != nil {
			return err
		}

		// 3️⃣ Delete game type data
		if err := repo.DeletePlayTypeData(id); err != nil {
			return err
		}

		// 3️⃣ Delete game type data
		// if err := repo.DeleteGameGameType(id); err != nil {
		//     return err
		// }

		if err := repo.Delete(id); err != nil {
			return err
		}

		return nil
	})

	if result == nil {
		tables, err := helper.GetTablesByIDs([]uint{}, 0, uint(id))
		if err == nil {
			BulkUpdateTableInfo(tables, "disable")
		}
	}

	return result
}

func (s *PlayService) GetPlayTypesByPlayID(id uint64) ([]models.GameType, error) {
	return s.Repo.FindPlayTypesByPlayID(id)
}
func (s *PlayService) AuditAction(
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
func (s *PlayService) FindByID(id uint64) (*models.Game, error) {
	return s.Repo.FindByID(id)
}
