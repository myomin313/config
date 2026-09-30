package controllers

import (
	"config-service/ApiResponse"
	"config-service/Repositories"
	"config-service/Services"
	"config-service/config"

	"github.com/gofiber/fiber/v2"
)

// GameResult represents the structure of a game result report
type PlayResult struct {
	TableNo   string `json:"table_no"`
	Game      string `json:"game"`
	ShoeNo    int    `json:"shoe_no"`
	RoundNo   int    `json:"round_no"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Result    string `json:"result"`
	Status    string `json:"status"`
	VideoURL  string `json:"video_url"`
	Remark    string `json:"remark"`
}

// GetGameResultReport returns a list of game results
func GetPlayResultReport(c *fiber.Ctx) error {
	service := Services.RoundService{
		Repo: &Repositories.RoundRepository{
			DB: config.DB,
			// EngineDB: config.GameEngineDB,
		},
	}

	results, err := service.GetPlayResults()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Convert to response format
	var response []PlayResult
	for _, r := range results {
		response = append(response, PlayResult{
			TableNo:   r.TableNo,
			Game:      r.GameName,
			ShoeNo:    int(r.ShoeNo),
			RoundNo:   int(r.RoundNo),
			StartTime: r.StartTime,
			EndTime:   r.EndTime,
			Result:    r.Result,
			Status:    r.Status,
			VideoURL:  r.VideoURL,
			Remark:    r.Remark,
		})
	}

	// Return response in same format as /games API
	return ApiResponse.SuccessWithPagination(c, 200, len(response), response, "Game result report retrieved successfully")
}
