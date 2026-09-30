package models


type GameType struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	GameTypeName string `gorm:"type:varchar(255);not null"`
	GameDataID uint
	
}


func init() {
	RegisterModel(&GameType{})
}