package controllers

import (
	"log"
	"strconv"
	"time"

	"config-service/config"
	"config-service/helper"
	"config-service/requests"
	"config-service/utils"

	"config-service/ApiResponse"
	"config-service/Repositories"
	"config-service/Services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ------------------------ Controller Struct ------------------------

type ReportController struct {
	Service *Services.ReportService
}

func NewReportController(db *gorm.DB) *ReportController {
	auditService := &helper.AuditService{
		//Redis: redisConn,
		Channel: config.RabbitChannel,
		Queue:   "audit_logs",
	}
	return &ReportController{
		Service: &Services.ReportService{
			Repo: &Repositories.ReportRepository{
				DB: db,
			},
			Audit: auditService,
		},
	}
}

// ------------------------ Helper ------------------------

// ------------------------ Handlers ------------------------

// GetPlatforms godoc
// @Summary      List game tables
// @Description  Get all game tables
// @Tags         Gametables
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /reports [get]
func (pc *ReportController) GetPlayResult(c *fiber.Ctx) error {
	req := new(requests.GameResultReportSearchRequest)
	if err := c.QueryParser(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusBadRequest, "unable_to_parse_request_body", fullMessages[0])
	}
	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		log.Println("err", fullMessages[0])
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	if req.ResultTime != "" {
		if req.StartDate == "" || req.EndDate == "" {
			return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "start_date and end_date are required when result_time is provided")
		}

		loc, err := time.LoadLocation("Asia/Phnom_Penh")
		if err != nil {
			return ApiResponse.Error(c, fiber.StatusInternalServerError, "failed to load timezone")
		}

		var startTime, endTime time.Time
		layoutDatetime := "2006-01-02 15:04:05"
		layoutDate := "2006-01-02"

		startTime, err = time.ParseInLocation(layoutDatetime, req.StartDate, loc)
		if err != nil {
			startTime, err = time.ParseInLocation(layoutDate, req.StartDate, loc)
			if err != nil {
				return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "start_date must be in format YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
			}
		}

		endTime, err = time.ParseInLocation(layoutDatetime, req.EndDate, loc)
		endDateWasDateOnly := false
		if err != nil {
			endTime, err = time.ParseInLocation(layoutDate, req.EndDate, loc)
			if err != nil {
				return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "end_date must be in format YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
			}
			endDateWasDateOnly = true
		}

		if endDateWasDateOnly {
			endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}

		if startTime.After(endTime) {
			return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "start_date must be before or equal to end_date")
		}

		if endTime.Sub(startTime) > 30*24*time.Hour {
			return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", "date range must be 30 days or less")
		}
		log.Printf("Parsed startTime: %s, endTime: %s\n", startTime.Format(layoutDatetime), endTime.Format(layoutDatetime))
		req.StartDate = startTime.Format(layoutDatetime)
		req.EndDate = endTime.Format(layoutDatetime)
	}

	platforms, err := pc.Service.GetPlayResult(c, req)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// resp := make([]Resources.GameTableResource, 0, len(platforms))
	// for i := range platforms {
	// 	resp = append(resp, mapGameTableToResponse(&platforms[i]))
	// }

	return ApiResponse.JsonSuccess(c, 200, platforms, "Rule URL retrieved successfully")
}

func (pc *ReportController) UpdateRoundRemark(c *fiber.Ctx) error {
	roundIDStr := c.Params("id")
	if roundIDStr == "" {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "round id is required")
	}

	var body struct {
		Remark string `json:"remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	if body.Remark == "" {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "remark is required")
	}

	roundID, err := strconv.ParseUint(roundIDStr, 10, 64)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid round id")
	}

	result, err := pc.Service.UpdateRoundRemark(c, roundID, body.Remark)
	log.Println(result, err)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	creator, _ := c.Locals("user").(jwt.MapClaims)
	roundResp, ok := result["roundResp"].(map[string]interface{})
	if !ok {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, "invalid response")
	}

	oldRemark := map[string]interface{}{
		"game_name": roundResp["gameName"],
		"table_no":  roundResp["tableNo"],
		"shoe_no":   roundResp["shoeNo"],
		"round_no":  roundResp["roundNo"],
		"remarks":   roundResp["originalRemark"],
	}

	newRemark := map[string]interface{}{
		"game_name": roundResp["gameName"],
		"table_no":  roundResp["tableNo"],
		"shoe_no":   roundResp["shoeNo"],
		"round_no":  roundResp["roundNo"],
		"remarks":   body.Remark,
	}
	pc.Service.AuditAction(
		c,
		"update remark",
		"Game Result",
		creator,
		oldRemark,
		newRemark,
	)

	return ApiResponse.Success(c, fiber.StatusOK, nil, "Round remark updated successfully")
}

func (pc *ReportController) UpdateRoundResult(c *fiber.Ctx) error {
	roundIDStr := c.Params("id")
	if roundIDStr == "" {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "round id is required")
	}

	roundID, err := strconv.ParseUint(roundIDStr, 10, 64)
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid round id")
	}

	req := new(requests.UpdateRoundResultRequest)
	if err := c.BodyParser(req); err != nil {
		return ApiResponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(req); err != nil {
		fullMessages := utils.ValidationErrorMessages(err, req)
		return ApiResponse.ValidationError(c, fiber.StatusUnprocessableEntity, "validation_error", fullMessages[0])
	}

	result, err := pc.Service.UpdateRoundResult(c, roundID, req)
	log.Println(result, err)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument, codes.FailedPrecondition:
				log.Println("validation error:", s.Message())
				return ApiResponse.Error(c, fiber.StatusBadRequest, s.Message())
			case codes.NotFound:
				log.Println("not found error:", s.Message())
				return ApiResponse.Error(c, fiber.StatusNotFound, s.Message())
			}
		}
		log.Println("unexpected error:", err)
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return ApiResponse.Success(c, fiber.StatusOK, nil, "Round result updated successfully")
}
