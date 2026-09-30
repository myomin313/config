package requests

type GameResultReportSearchRequest struct {
	GameId     uint64 `query:"game_id" validate:"omitempty,required,min=1" example:"1,2,3"`
	ShoeNo     uint64 `query:"shoe_no" validate:"omitempty,required,min=1" example:"1,2,3"`
	RoundNo    uint64 `query:"round_no" validate:"omitempty,required,min=1" example:"1,2,3"`
	TableNo    string `query:"table_no" validate:"omitempty,required" example:"1,2,3"`
	ResultTime string `query:"result_time" validate:"omitempty,required,oneof=start_time end_time" example:"start_time"`
	StartDate  string `query:"start_date" validate:"omitempty,required" example:"2024-01-01 00:00:00"`
	EndDate    string `query:"end_date" validate:"omitempty,required" example:"2024-01-31"`
	Status     string `query:"status" validate:"omitempty,required,oneof=start_betting opening settled" example:"settled"`
}

type UpdateRoundResultRequest struct {
	Dice1  uint32 `json:"dice1" validate:"required,gte=1,lte=6"`
	Dice2  uint32 `json:"dice2" validate:"required,gte=1,lte=6"`
	Dice3  uint32 `json:"dice3" validate:"required,gte=1,lte=6"`
	Remark string `json:"remark" validate:"required"`
}
