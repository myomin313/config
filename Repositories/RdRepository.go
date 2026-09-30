package Repositories

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// RoundRepository handles DB operations for game results
type RoundRepository struct {
	DB       *gorm.DB // Admin service DB
	EngineDB *gorm.DB // Game engine service DB
}

// GameResultRow represents the joined data for game result report
type PlayResultRow struct {
	TableNo   string `gorm:"column:table_number"`
	GameName  string `gorm:"column:game_name"`
	ShoeNo    int64  `gorm:"column:shoe_no"`
	RoundNo   int64  `gorm:"column:round_no"`
	StartTime string `gorm:"column:start_time"`
	EndTime   string `gorm:"column:end_time"`
	Result    string `gorm:"column:result"`
	Status    string `gorm:"column:status"`
	VideoURL  string `gorm:"column:video_url"`
	Remark    string `gorm:"column:remark"`
}

// RoundData represents a single round from game engine
type RoundData struct {
	ID             int64  `gorm:"column:id"`
	TableID        int64  `gorm:"column:table_id"`
	ShoeNo         int64  `gorm:"column:shoe_no"`
	RoundNo        int64  `gorm:"column:round_no"`
	Status         string `gorm:"column:status"`
	Dice1          int    `gorm:"column:dice1"`
	Dice2          int    `gorm:"column:dice2"`
	Dice3          int    `gorm:"column:dice3"`
	BettingStartAt string `gorm:"column:betting_start_at"`
	SettledAt      string `gorm:"column:settled_at"`
}

// GameTableRow represents game table data
type PlayTableRow struct {
	ID        int64  `gorm:"column:id"`
	TableNo   string `gorm:"column:table_number"`
	GameType  int64  `gorm:"column:game_type"`
	Status    string `gorm:"column:status"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
}

// GameDataRow represents game data
type PlayDataRow struct {
	ID   int64  `gorm:"column:id"`
	Name string `gorm:"column:name_en"`
}

// GetGameResults returns game result data combining admin and engine service data
func (r *RoundRepository) GetPlayResults(results *[]PlayResultRow) error {
	// If game engine DB is not available, return only admin service data
	if r.EngineDB == nil {
		log.Println("⚠️  Game Engine DB not available, returning admin service data only")
		return r.getPlayResultsFromAdminOnly(results)
	}

	// Query rounds from game engine database
	var rounds []RoundData
	query := `
		SELECT
			id,
			table_id,
			shoe_no,
			round_no,
			status,
			dice1,
			dice2,
			dice3,
			betting_start_at::text,
			settled_at::text
		FROM rounds
		ORDER BY created_at DESC
		LIMIT 100
	`

	if err := r.EngineDB.Raw(query).Scan(&rounds).Error; err != nil {
		log.Printf("❌ Error querying rounds from game engine: %v\n", err)
		return err
	}

	if len(rounds) == 0 {
		log.Println("⚠️  No rounds found in game engine database")
		return nil
	}

	log.Printf("✅ Found %d rounds from game engine\n", len(rounds))

	// Query game tables from admin database
	playTableMap := make(map[int64]*PlayTableRow)

	playTableQuery := `
		SELECT
			id,
			table_number,
			game_type,
			status,
			created_at::text,
			updated_at::text
		FROM game_tables
	`

	var playTableRows []PlayTableRow
	if err := r.DB.Raw(playTableQuery).Scan(&playTableRows).Error; err != nil {
		log.Printf("❌ Error querying game tables: %v\n", err)
	}

	for i := range playTableRows {
		playTableMap[playTableRows[i].ID] = &playTableRows[i]
	}

	log.Printf("✅ Found %d game tables\n", len(playTableRows))

	// Query game data from admin database
	playDataMap := make(map[int64]string)

	playDataQuery := `
		SELECT
			g.id,
			gd.name_en
		FROM games g
		LEFT JOIN game_data gd ON g.game_data_id = gd.id
	`

	var playDataRows []PlayDataRow
	if err := r.DB.Raw(playDataQuery).Scan(&playDataRows).Error; err != nil {
		log.Printf("❌ Error querying game data: %v\n", err)
	}

	for _, gd := range playDataRows {
		playDataMap[gd.ID] = gd.Name
	}

	log.Printf("✅ Found %d games\n", len(playDataRows))

	// Combine data
	for _, round := range rounds {
		playTable, exists := playTableMap[round.TableID]
		if !exists {
			log.Printf("⚠️  Game table not found for table_id: %d\n", round.TableID)
			continue
		}

		playName := playDataMap[playTable.GameType]
		if playName == "" {
			playName = "Unknown Game"
		}

		result := PlayResultRow{
			TableNo:   playTable.TableNo,
			GameName:  playName,
			ShoeNo:    round.ShoeNo, // Add shoe_no if available in game engine
			RoundNo:   round.RoundNo,
			StartTime: round.BettingStartAt,
			EndTime:   round.SettledAt,
			Result:    fmt.Sprintf("%d,%d,%d", round.Dice1, round.Dice2, round.Dice3),
			Status:    round.Status,
			VideoURL:  "",
			Remark:    "",
		}

		*results = append(*results, result)
	}

	log.Printf("✅ Combined %d game results\n", len(*results))
	return nil
}

func (r *RoundRepository) getPlayResultsFromAdminOnly(results *[]PlayResultRow) error {
	query := `
		SELECT
			gt.table_number::text,
			gd.name_en::text,
			0::bigint as shoe_no,
			0::bigint as round_no,
			gt.created_at::text as start_time,
			gt.updated_at::text as end_time,
			'' as result,
			gt.status::text,
			'' as video_url,
			'' as remark
		FROM game_tables gt
		LEFT JOIN games g ON gt.game_type = g.id
		LEFT JOIN game_data gd ON g.game_data_id = gd.id
		ORDER BY gt.created_at DESC
		LIMIT 100
	`
	return r.DB.Raw(query).Scan(results).Error
}
