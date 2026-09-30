package Services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"config-service/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PlayTableService struct {
	Repo           *Repositories.PlayTableRepository
	PlayRepo       *Repositories.PlayRepository
	RedisPublisher *RedisPublisher
	Audit          *helper.AuditService
}

var t = time.Now()

// Create a new Platform
func (s *PlayTableService) Create(req *requests.GameTableCreateRequest, id uint) (*models.GameTable, error) {

	tableNumber := strings.TrimSpace(req.TableNumber)
	exists, err := s.Repo.TableNumberExists(tableNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("table number already exists")
	}
	playtable := &models.GameTable{
		AffiliateOffice: req.AffiliateOffice,
		GameType:        req.GameType,
		GameTypeID:      req.GameTypeID,
		BettingTime:     req.BettingTime,
		SettlementTime:  req.SettlementTime,
		TableNumber:     tableNumber,
		Status:          req.Status,
		//TableResultType: req.TableResultType,
		MobileIcon:  req.MobileIcon,
		Sort:        req.Sort,
		PCIcon:      req.PCIcon,
		Mode:        req.Mode,
		Topic:       req.Topic,
		CreatedByID: id,
		UpdatedByID: id,
	}

	// Game translations
	for lang, name := range req.GameTableName {
		playtable.Translations = append(playtable.Translations, models.GameTranslation{
			LanguageCode: lang,
			LangName:     "game_table_name",
			Content:      name,
			Type:         "api",
			TableType:    "game_tables",
		})
	}
	// resolution data
	for _, resolution := range req.Resolutions {

		playtable.Resolutions = append(
			playtable.Resolutions,
			models.GameTableResolution{
				Resolution: resolution.Resolution,
				RoomName:   strings.TrimSpace(resolution.RoomName),
			},
		)
	}

	if err := s.Repo.Create(playtable); err != nil {
		return nil, err
	}

	return playtable, nil
}

// All returns all Platforms
func (s *PlayTableService) All(req *requests.GameTableSearchRequest, lang string) ([]models.GameTable, error) {
	var platforms []models.GameTable

	if err := s.Repo.FindAll(req, &platforms, lang); err != nil {
		return nil, err
	}
	return platforms, nil
}

// FindByID returns a Platform by ID
func (s *PlayTableService) FindByID(id uint64) (*models.GameTable, error) {
	var platform models.GameTable
	if err := s.Repo.FindByID(id, &platform); err != nil {
		return nil, errors.New("Record not found")
	}
	return &platform, nil
}

// Update a Platform by ID
// func (s *GameTableService) Update(id uint64, req *requests.GameTableUpdateRequest, userId uint) (*models.GameTable, error) {
// 	tableNumber := strings.TrimSpace(*req.TableNumber)
// 	if req.TableNumber != nil {
// 		exists, err := s.Repo.TableNumberExistsExcludeID(tableNumber, id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if exists {
// 			return nil, errors.New("table number already exists")
// 		}
// 	}

// 	platform, err := s.FindByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 🔥 Apply only non-nil fields from request to model
// 	basicReq := requests.GameTableBasicUpdateRequest{
// 		GameTableName:   req.GameTableName,
// 		AffiliateOffice: req.AffiliateOffice,
// 		GameType:        req.GameType,
// 		GameTypeID:      req.GameTypeID,
// 		BettingTime:     req.BettingTime,
// 		SettlementTime:  req.SettlementTime,
// 		TableNumber:     req.TableNumber,
// 		Status:          req.Status,
// 		TableResultType: req.TableResultType,
// 		MobileIcon:      req.MobileIcon,
// 		PCIcon:          req.PCIcon,
// 		Sort:            req.Sort,
// 		Mode:            req.Mode,
// 		Topic:           req.Topic,
// 	}

// 	utils.UpdateStruct(platform, &basicReq)
// 	//utils.UpdateStruct(platform, req)

// 	platform.UpdatedByID = userId
// 	platform.UpdatedAt = time.Now()
// 	platform.TableNumber = tableNumber
// 	if req.GameTypeID != nil {
// 		platform.GameTypeID = *req.GameTypeID
// 	}

// 	if req.GameTableName != nil && len(*req.GameTableName) > 0 {
// 		// ---- Replace translations ----
// 		if errMsg := s.GameRepo.DeleteGameTranslations(id, "game_tables"); errMsg != nil {
// 			return nil, err
// 		}
// 		platform.Translations = nil
// 		for lang, name := range *req.GameTableName {
// 			platform.Translations = append(platform.Translations, models.GameTranslation{
// 				LanguageCode: lang,
// 				LangName:     "game_table_name",
// 				Content:      name,
// 				Type:         "api",
// 				TableType:    "game_tables",
// 			})
// 		}
// 	}
// 	log.Println("platform AffiliateOffice", platform.AffiliateOffice)

// 	if req.Resolutions != nil {

// 		// Delete old resolutions
// 		if err := s.Repo.DeleteResolutions(id); err != nil {
// 			return nil, err
// 		}

// 		// Clear existing association in memory
// 		platform.Resolutions = nil

// 		// Convert REQUEST models -> DATABASE models
// 		for _, resolution := range *req.Resolutions {

// 			platform.Resolutions = append(
// 				platform.Resolutions,
// 				models.GameTableResolution{
// 					GameTableID: uint(id),
// 					Resolution:  resolution.Resolution,
// 					RoomName:    strings.TrimSpace(resolution.RoomName),
// 				},
// 			)
// 		}
// 	}
// 	if errRRR := s.Repo.Update(platform); errRRR != nil {
// 		log.Println("game table update error", errRRR)
// 		return nil, err
// 	}

//		return platform, nil
//	}
func (s *PlayTableService) Update(
	id uint64,
	req *requests.GameTableUpdateRequest,
	userId uint,
) (*models.GameTable, error) {

	// =========================================================
	// 1. Validate table number
	// =========================================================

	var tableNumber string

	if req.TableNumber != nil {
		tableNumber = strings.TrimSpace(*req.TableNumber)

		exists, err := s.Repo.TableNumberExistsExcludeID(tableNumber, id)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errors.New("table number already exists")
		}
	}

	// =========================================================
	// 2. Get existing game table
	// =========================================================

	platform, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	// =========================================================
	// 3. Apply basic fields
	// =========================================================

	basicReq := requests.GameTableBasicUpdateRequest{
		GameTableName:   req.GameTableName,
		AffiliateOffice: req.AffiliateOffice,
		GameType:        req.GameType,
		GameTypeID:      req.GameTypeID,
		BettingTime:     req.BettingTime,
		SettlementTime:  req.SettlementTime,
		TableNumber:     req.TableNumber,
		Status:          req.Status,
		//TableResultType: req.TableResultType,
		MobileIcon: req.MobileIcon,
		PCIcon:     req.PCIcon,
		Sort:       req.Sort,
		Mode:       req.Mode,
		Topic:      req.Topic,
	}

	utils.UpdateStruct(platform, &basicReq)

	platform.UpdatedByID = userId
	platform.UpdatedAt = time.Now()

	if req.TableNumber != nil {
		platform.TableNumber = tableNumber
	}

	if req.GameTypeID != nil {
		platform.GameTypeID = *req.GameTypeID
	}

	// =========================================================
	// 4. Update translations
	// =========================================================

	if req.GameTableName != nil && len(*req.GameTableName) > 0 {

		if errMsg := s.PlayRepo.DeletePlayTranslations(
			id,
			"game_tables",
		); errMsg != nil {
			return nil, errMsg
		}

		platform.Translations = nil

		for lang, name := range *req.GameTableName {
			platform.Translations = append(
				platform.Translations,
				models.GameTranslation{
					LanguageCode: lang,
					LangName:     "game_table_name",
					Content:      name,
					Type:         "api",
					TableType:    "game_tables",
				},
			)
		}
	}

	// =========================================================
	// 5. IMPORTANT
	//
	// Repository uses:
	//
	// FullSaveAssociations: true
	//
	// Therefore don't allow GORM to manage resolutions here.
	// We will manage resolutions manually below.
	// =========================================================

	platform.Resolutions = nil

	// =========================================================
	// 6. Update main GameTable
	// =========================================================

	if err := s.Repo.Update(platform); err != nil {
		log.Println("game table update error:", err)
		return nil, err
	}

	// =========================================================
	// 7. Resolution handling
	// =========================================================
	//
	// Rules:
	//
	// mode = digital
	//     -> delete all resolutions
	//
	// mode = dealer
	//     + resolutions provided
	//         -> delete old + create new
	//
	// mode = dealer
	//     + resolutions omitted
	//         -> keep existing resolutions
	//
	// mode = dealer
	//     + resolutions: []
	//         -> delete all resolutions
	// =========================================================

	if req.Mode != nil {

		mode := strings.TrimSpace(*req.Mode)

		log.Println("Game table mode:", mode)

		// -----------------------------------------------------
		// Digital mode
		// -----------------------------------------------------

		if mode != "dealer" {

			log.Println(
				"Mode is not dealer, deleting resolutions for table:",
				id,
			)

			if err := s.Repo.DeleteResolutions(id); err != nil {
				return nil, fmt.Errorf(
					"failed to delete resolutions: %w",
					err,
				)
			}
		}

		// -----------------------------------------------------
		// Dealer mode
		// -----------------------------------------------------

		if mode == "dealer" && req.Resolutions != nil {

			log.Println(
				"Updating resolutions for dealer table:",
				id,
			)

			// Delete old resolutions first
			if err := s.Repo.DeleteResolutions(id); err != nil {
				return nil, fmt.Errorf(
					"failed to delete old resolutions: %w",
					err,
				)
			}

			// Create new resolutions
			for _, resolution := range *req.Resolutions {

				newResolution := models.GameTableResolution{
					GameTableID: uint(id),
					Resolution:  resolution.Resolution,
					RoomName:    strings.TrimSpace(resolution.RoomName),
				}

				if err := s.Repo.CreateResolution(&newResolution); err != nil {
					return nil, fmt.Errorf(
						"failed to create resolution: %w",
						err,
					)
				}
			}
		}
	}

	// =========================================================
	// 8. Reload complete GameTable
	// =========================================================

	platform, err = s.FindByID(id)
	if err != nil {
		return nil, err
	}

	return platform, nil
}

// Delete a Platform by ID
func (s *PlayTableService) Delete(id uint64) error {
	platform, err := s.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.PlayRepo.DeletePlayTranslations(id, "game_tables"); err != nil {
		return fmt.Errorf("failed to delete game table translations: %w", err)
	}
	// Delete resolutions
	if err := s.Repo.DeleteResolutions(id); err != nil {
		return fmt.Errorf(
			"failed to delete game table resolutions: %w",
			err,
		)
	}
	if err := s.Repo.Delete(platform); err != nil {
		return fmt.Errorf("failed to delete game table: %w", err)
	}
	s.RedisPublisher.DisableTableInfo(platform)
	return nil
}

// update the gametable status
func (s *PlayTableService) MassUpdateStatus(
	tableIDs []uint,
	status string,
) error {
	resp := s.Repo.MassUpdateStatus(tableIDs, status)

	if resp == nil {
		tables, err := helper.GetTablesByIDs(tableIDs, 0, 0)
		if err == nil {
			BulkUpdateTableInfo(tables, status)
		}
	}

	return resp
}

// All returns all Platforms
func (s *PlayTableService) AllForFE(req *requests.GameTableFESearchRequest, lang string) ([]models.GameTable, error) {
	var platforms []models.GameTable

	if err := s.Repo.FindAllForFE(req, &platforms, lang); err != nil {
		return nil, err
	}
	return platforms, nil
}
func (s *PlayTableService) AuditAction(
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
