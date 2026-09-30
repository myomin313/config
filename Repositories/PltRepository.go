package Repositories

import (
	"config-service/models"
	"config-service/requests"
	"fmt"

	"log"

	"gorm.io/gorm"
)

// GameTableRepository handles DB operations for Platform
type PlayTableRepository struct {
	DB *gorm.DB
}

// Create inserts a new Platform into the database
func (r *PlayTableRepository) Create(playTable *models.GameTable) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(playTable).Error
}

// FindAll returns all Platforms
func (r *PlayTableRepository) FindAll(req *requests.GameTableSearchRequest, playTables *[]models.GameTable, lang string) error {
	query := r.DB.Preload("Translations", "lang_name = ?", "game_table_name").Preload("Resolutions")

	// 🔍 Table name search
	// 🔍 Translation table search (FIXED)
	if req.TableName != nil {
		like := "%" + *req.TableName + "%"

		query = query.
			Joins(`
				JOIN game_translations t
				ON t.table_id = game_tables.id
			`).
			Where("t.lang_name = ?", "game_table_name").
			Where("t.language_code = ?", lang).
			Where(`
				t.content LIKE ?
				`, like).
			Group("game_tables.id") // avoid duplicates
	}

	// 🎮 Game type filter
	if req.GameType != nil {
		query = query.Where("game_type = ?", *req.GameType)
	}

	// 📌 Status filter
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	query = query.Order("game_tables.sort ASC")
	return query.Find(playTables).Error

}

// FindByID returns a Platform by ID
func (r *PlayTableRepository) FindByID(id uint64, playTable *models.GameTable) error {
	return r.DB.
		Preload("Translations", "table_id = ? AND lang_name = ?", id, "game_table_name").
		Preload("Resolutions").
		First(playTable, id).Error
}

// Update updates an existing Platform
func (r *PlayTableRepository) Update(playTable *models.GameTable) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(playTable).Error
}

// Delete deletes a Platform
func (r *PlayTableRepository) Delete(playTable *models.GameTable) error {
	return r.DB.Delete(playTable).Error
}

func (r *PlayTableRepository) MassUpdateStatus(
	tableIDs []uint,
	status string,
) error {

	result := r.DB.
		Model(&models.GameTable{}).
		Where("id IN ?", tableIDs).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no records updated")
	}

	return nil
}

func (r *PlayTableRepository) GetAllTables() ([]models.GameTable, error) {
	var tables []models.GameTable // Change to a slice

	// Use .Find() to fetch all matching records
	if err := r.DB.Find(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}

// FindAll returns all Platforms
func (r *PlayTableRepository) FindAllForFE(req *requests.GameTableFESearchRequest, playTables *[]models.GameTable, lang string) error {
	query := r.DB.
		Preload("Translations", "lang_name = ? AND language_code = ?", "game_table_name", lang)
	// Affiliate Office filter
	if req.AffiliateOffice != nil {
		query = query.Where("game_tables.affiliate_office = ?", *req.AffiliateOffice)
	}
	// Game type filter
	if req.Game != nil {
		query = query.Where("game_tables.game_type = ?", *req.Game)
	}
	// Always ENABLE only
	query = query.Where("game_tables.status = ?", "active")
	// 🎯 Sorting
	query = query.Order("game_tables.sort ASC")
	return query.Find(playTables).Error
}

func (r *PlayTableRepository) TableNumberExists(TableNumber string) (bool, error) {
	log.Println("🔍 Checking table number:", TableNumber)

	var count int64
	err := r.DB.
		Model(&models.GameTable{}). // ✅ IMPORTANT
		Where("table_number = ?", TableNumber).
		Count(&count).Error

	if err != nil {
		log.Println("❌ DB error:", err)
		return false, err
	}

	log.Println("✅ Count result:", count)

	return count > 0, nil
}
func (r *PlayTableRepository) TableNumberExistsExcludeID(TableNumber string, id uint64) (bool, error) {
	var count int64

	err := r.DB.Model(&models.GameTable{}).
		Where("table_number = ? AND id <> ?", TableNumber, id).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// get tables by id
// by hall_id or game_id or table_ids
func (r *PlayTableRepository) GetTablesByIDs(tableIDs []uint, hallID uint, gameId uint) ([]models.GameTable, error) {
	var tables []models.GameTable
	sql := r.DB.
		Select("id, table_number, affiliate_office, game_type, betting_time, settlement_time, status, mode, topic")

	if len(tableIDs) > 0 {
		sql = sql.Where("id IN ?", tableIDs)
	}

	if hallID > 0 {
		sql = sql.Where("affiliate_office = ?", hallID)
	}

	if gameId > 0 {
		sql = sql.Where("game_type = ?", gameId)
	}

	result := sql.Find(&tables)

	if result.Error != nil {
		return nil, result.Error
	}

	return tables, nil
}
func (r *PlayTableRepository) DeleteResolutions(
	playTableID uint64,
) error {
	return r.DB.
		Where("game_table_id = ?", playTableID).
		Delete(&models.GameTableResolution{}).
		Error
}
func (r *PlayTableRepository) CreateResolution(
	resolution *models.GameTableResolution,
) error {
	return r.DB.Create(resolution).Error
}
