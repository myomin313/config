package migrations

import "config-service/migrations/seeders"

func Seeders() {
	seeders.SeedGameTypes()
	seeders.CreateMasterGameTypes()
	seeders.CreatePlatforms()
	seeders.SeedGameData()
	seeders.Maintenance()
	seeders.SeedMerchantCurrency()
	//seeders.SeedCurrency()
}
