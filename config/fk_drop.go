package config

import (
	"fmt"
	"log"
)

// func RunRawMigrations() {
// 	sql := `
// 	DO $$
// 	BEGIN
// 		IF EXISTS (
// 			SELECT 1
// 			FROM information_schema.table_constraints
// 			WHERE constraint_name = 'fk_games_translations'
// 			AND table_name = 'game_translations'
// 		) THEN
// 			ALTER TABLE game_translations
// 			DROP CONSTRAINT fk_games_translations;
// 		END IF;
// 	END $$;
// 	`

// 	if err := DB.Exec(sql).Error; err != nil {
// 		log.Fatalf("❌ Raw migration failed: %v", err)
// 	}

// 	log.Println("✅ Raw migrations executed")
// }

// RunRawMigrations checks if the foreign key exists and drops it
func RunRawMigrations(tableName string, fkName string) {

	var exists bool

	// Check if the foreign key exists
	checkSQL := `
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints
            WHERE table_name = ? 
              AND constraint_name = ? 
              AND constraint_type = 'FOREIGN KEY'
        )
    `
	if err := DB.Raw(checkSQL, tableName, fkName).Scan(&exists).Error; err != nil {
		log.Fatalf("❌ Failed to check foreign key existence: %v", err)
	}

	// Drop the foreign key if it exists
	if exists {
		dropSQL := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", tableName, fkName)
		if err := DB.Exec(dropSQL).Error; err != nil {
			log.Fatalf("❌ Failed to drop foreign key: %v", err)
		}
		log.Println("✅ Foreign key dropped:", fkName)
	} else {
		log.Println("ℹ️ Foreign key does not exist, skipping:", fkName)
	}

	log.Println("✅ Raw migrations executed")
}
