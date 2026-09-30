package Repositories

import (
	"config-service/models"
	"sort"
	"strings"

	//"log"

	"gorm.io/gorm"
)

type PlayRepository struct {
	DB *gorm.DB
}

// Fetch all games with related data
// func (r *GameRepository) FindAll() ([]models.Game, error) {
// 	var games []models.Game
// 	err := r.DB.
// 		//Preload("CreatedByUser").
// 		//Preload("UpdatedByUser").
// 		Preload("GameTypes").
// 		Preload("Translations", "lang_name = ? ", "game_name").
// 		Preload("GameTypeData.MasterGameType.GameTypeTranslation"). // preload master for names

// 		Find(&games).Error

// 	return games, err
// }

// Fetch all games with related data
func (r *PlayRepository) FindAllEnable() ([]models.Game, error) {
	var plays []models.Game

	err := r.DB.
		//Preload("CreatedByUser").
		//Preload("UpdatedByUser").
		Where("status = ? ", "enable").
		Preload("GameTypes").
		Preload("Translations", "lang_name = ? ", "game_name").
		Preload("GameTypeData.MasterGameType.GameTypeTranslation"). // preload master for names
		Find(&plays).Error

	return plays, err
}

func (r *PlayRepository) FindPlayTypesByPlayID(playID uint64) ([]models.GameType, error) {
	var playTypes []models.GameType

	err := r.DB.
		Table("game_types").
		Select("game_types.id, game_types.game_type_name").
		Joins("JOIN game_game_types ggt ON ggt.game_type_id = game_types.id").
		Where("ggt.game_id = ?", playID).
		Scan(&playTypes).Error

	if err != nil {
		return nil, err
	}

	return playTypes, nil
}

// Create a game with related GameTypeData and translations
func (r *PlayRepository) Create(play *models.Game) error {
	return r.DB.Create(play).Error
}

// Find a game by ID with all relations
func (r *PlayRepository) FindByID(id uint64) (*models.Game, error) {
	var play models.Game

	if err := r.DB.
		Preload("GameTypes").Find(&play, id).
		Preload("GameTypeData.MasterGameType.GameTypeTranslation").
		//Preload("CreatedByUser").
		//Preload("UpdatedByUser").
		//Preload("Translations").
		Preload("Translations", "table_id = ? AND lang_name = ?", id, "game_name").
		First(&play, id).Error; err != nil {

		return nil, err
	}
	//log.Println(game)
	return &play, nil
}

func (r *PlayRepository) PlayTypeExists(id uint64) (bool, error) {

	var count int64
	err := r.DB.
		Model(&models.GameType{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}
func (r *PlayRepository) MasterPlayTypeExists(id uint64) (bool, error) {
	var count int64
	err := r.DB.
		Model(&models.MasterGameType{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

// Delete all translations for a game
func (r *PlayRepository) DeletePlayTranslations(playID uint64, label string) error {
	return r.DB.
		Where("table_id = ? AND table_type = ?", playID, label).
		Delete(&models.GameTranslation{}).Error
}

// func (r *GameRepository) DeleteGameTranslations(gameID uint64,  tableType string) error {
// 	return r.DB.
// 		Where("table_id = ? AND table_type", gameID, tableType).
// 		Delete(&models.GameTranslation{}).Error
// }

// Delete all GameTypeData for a game
func (r *PlayRepository) DeletePlayTypeData(playID uint64) error {
	return r.DB.
		Where("game_id = ?", playID).
		Delete(&models.GameTypeData{}).Error
}

// Update game with associations (translations + GameTypeData)
//
//	func (r *GameRepository) Update(game *models.Game) error {
//		return r.DB.
//			Session(&gorm.Session{FullSaveAssociations: true}).
//			Save(game).Error
//	}
func (r *PlayRepository) Update(play *models.Game) error {
	return r.DB.Save(play).Error // do NOT use FullSaveAssociations
}

// Delete a game by ID (will not automatically delete translations/GameTypeData if cascade is not configured)
func (r *PlayRepository) Delete(playID uint64) error {
	// Optional: delete related data first
	// r.DeleteGameTranslations(gameID)
	// r.DeleteGameTypeData(gameID)

	return r.DB.Unscoped().Delete(&models.Game{}, playID).Error
}

func (r *PlayRepository) ReplacePlayTypes(play *models.Game, playTypes []models.GameType) error {
	return r.DB.Model(play).Association("GameTypes").Replace(playTypes)
}

func (r *PlayRepository) FindByPlayDataID(playDataID uint) (*models.Game, error) {
	var play models.Game
	err := r.DB.Where("game_data_id = ?", playDataID).First(&play).Error
	if err != nil {
		return nil, err
	}
	return &play, nil
}
func (r *PlayRepository) FindAll() ([]models.Game, error) {
	var plays []models.Game
	err := r.DB.
		Preload("GameTypes").
		Preload("Translations", "lang_name = ? ", "game_name").
		Preload("GameTypeData.MasterGameType.GameTypeTranslation").
		Find(&plays).Error
	if err != nil {
		return nil, err
	}

	// sort GameTypeData per game by play order
	for gi := range plays {
		sort.SliceStable(plays[gi].GameTypeData, func(i, j int) bool {
			return playOrderIndex(plays[gi].GameTypeData[i].MasterGameType) <
				playOrderIndex(plays[gi].GameTypeData[j].MasterGameType)
		})
	}

	return plays, nil
}

func playOrderIndex(m *models.MasterGameType) int {
	if m == nil {
		return 999
	}
	order := map[string]int{
		"tiger":         1,
		"gold_coin":     2,
		"klouk":         3,
		"chicken":       4,
		"shrimp":        5,
		"fish":          6,
		"crab":          7,
		"triple_symbol": 8,
		"double_animal": 9,
		"triple_color":  10,
	}
	if v, ok := order[strings.TrimSpace(strings.ToLower(m.KeyName))]; ok {
		return v
	}
	return 999
}
