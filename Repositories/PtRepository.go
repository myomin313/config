package Repositories

import (
	"config-service/models"
	"config-service/requests"
	"log"

	"gorm.io/gorm"
)

// GameRedLimitRepository handles DB operations for Platform
type PlayRedLimitRepository struct {
	DB *gorm.DB
}

// Create inserts a new Platform into the database
func (r *PlayRedLimitRepository) Create(redlimit *models.GameRedLimit) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(redlimit).Error
}

// FindAll returns all Platforms
func (r *PlayRedLimitRepository) FindAll(req *requests.RedLimitSearchRequest,
	playTables *[]models.GameRedLimit,
	lang string,
) error {

	query := r.DB.
		Model(&models.GameRedLimit{}).
		Preload("Translations", "lang_name = ?", "red_limit").
		Preload("RedLimitTable", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		})

	if req.CurrencyID != nil {
		query = query.Where(
			r.DB.
				Where("game_red_limits.game_type_currency = ?", *req.CurrencyID).
				Or("game_red_limits.game_play_type_currency = ?", *req.CurrencyID),
		)
	}

	// 🔍 Translation table search (FIXED)
	if req.TableName != nil {
		content := "%" + *req.TableName + "%"
		log.Print("content", content)
		log.Print("lang", lang)
		query = query.
			Joins(`
			JOIN game_translations t
			ON t.table_id = game_red_limits.id
			AND t.table_type = ?
		`, "game_red_limits").
			Where("t.lang_name = ?", "red_limit").
			Where("t.language_code = ?", lang).
			Where(`
				t.content LIKE ?
				`, content).
			Group("game_red_limits.id") // avoid duplicates
	}
	// ✅ Custom sort column
	query = query.Order("game_red_limits.sort ASC")
	return query.Find(playTables).Error
}

func (r *PlayRedLimitRepository) FindByIDForDelete(id uint64, playTable *models.GameRedLimit) error {
	return r.DB.First(playTable, id).Error
}

// FindByID returns a Platform by ID
func (r *PlayRedLimitRepository) FindByID(id uint64, playRedLimit *models.GameRedLimit) error {
	// 🌐 Conditionally preload translations
	// if lang == nil || *lang == "" {
	// 	log.Println("I am here")
	// 	return r.DB.First(gameRedLimit, id).Error
	// }

	return r.DB.
		Preload("Translations", "table_id = ? AND lang_name = ?", id, "red_limit").
		Preload("RedLimitTable", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC") // keep base order stable
		}).
		First(playRedLimit, uint(id)).Error
}

// Update updates an existing Platform
func (r *PlayRedLimitRepository) Update(playTable *models.GameRedLimit) error {

	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(playTable).Error
}

// Delete deletes a Platform
func (r *PlayRedLimitRepository) Delete(playTable *models.GameRedLimit) error {
	return r.DB.Delete(playTable).Error
}

// Delete deletes a Platform
func (r *PlayRedLimitRepository) DeleteRedlimitTable(id uint64) error {
	return r.DB.
		Where("red_limit_id = ?", id).
		Delete(&models.RedLimitTable{}).
		Error
}
