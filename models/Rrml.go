package models

import "time"

type RoundResultModificationLog struct {
	ID                        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoundID                   uint64    `gorm:"index" json:"round_id"`
	TableNo                   string    `json:"table_no"`
	ShoeNo                    uint64    `json:"shoe_no"`
	RoundNo                   uint64    `json:"round_no"`
	GameID                    uint64    `json:"game_id"`
	InitialOriginalResult     string    `gorm:"type:text" json:"initial_original_result"`
	CurrentResultBeforeModify string    `gorm:"type:text" json:"current_result_before_modify"`
	ModifiedResult            string    `gorm:"type:text" json:"modified_result"`
	ModificationReason        string    `gorm:"type:text" json:"modification_reason"`
	OperatorAccount           string    `json:"operator_account"`
	OperationTime             time.Time `json:"operation_time"`
	OperationIP               string    `json:"operation_ip"`
	OperationDevice           string    `json:"operation_device"`
	OperationResult           string    `json:"operation_result"`
	CreatedAt                 time.Time `json:"created_at"`
}

func (RoundResultModificationLog) TableName() string {
	return "round_result_modification_logs"
}

func init() {
	RegisterModel(&RoundResultModificationLog{})
}
