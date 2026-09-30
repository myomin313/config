package Services

import (
	"config-service/Repositories"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/pb/commonpb"
	"config-service/requests"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ReportService struct {
	Repo  *Repositories.ReportRepository
	Audit *helper.AuditService
}

func NewReportService(repo *Repositories.ReportRepository) *ReportService {
	return &ReportService{Repo: repo}
}

// GetAll returns all banners
func (s *ReportService) GetPlayResult(c *fiber.Ctx, req *requests.GameResultReportSearchRequest) (map[string]interface{}, error) {
	pagination := helper.GetPagination(c)
	// 2. Fetch from repo
	//game, total, err := s.Repo.GetGameResult(pagination.Offset, pagination.Limit)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("Fetched %d game records, total: %d\n", len(game), total)
	conn := config.GameEngineGrpcClient()
	defer conn.Close()

	client := commonpb.NewRoundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Requesting game results with filters - GameId: %d, ShoeNo: %d, RoundNo: %d, TableNo: %s, Status: %s, StartDate: %s, EndDate: %s, ResultTime: %s\n", req.GameId, req.ShoeNo, req.RoundNo, req.TableNo, req.Status, req.StartDate, req.EndDate, req.ResultTime)
	//resp, err := client.GetRound(ctx, &commonpb.RoundEmpty{})
	resp, err := client.GetRound(ctx, &commonpb.RoundEmpty{
		PageNo:     int32(pagination.PageNo),
		PageSize:   int32(pagination.PageSize),
		ShoeNo:     req.ShoeNo,
		RoundNo:    req.RoundNo,
		TableNo:    req.TableNo,
		Status:     req.Status,
		GameId:     req.GameId,
		EndDate:    req.EndDate,
		StartDate:  req.StartDate,
		ResultTime: req.ResultTime,
	})
	// if err != nil {
	// 	log.Println("Error calling GetRound:", err)
	// }
	// log.Printf("gRPC Response: %+v\n", resp)

	if err != nil || resp == nil {
		log.Println("Error calling GetRound:", err)

		return map[string]interface{}{
			"records":    []interface{}{},
			"page_no":    pagination.PageNo,
			"page_size":  pagination.PageSize,
			"total":      0,
			"total_page": 0,
		}, nil
	}

	acceptLanguage := c.Get("Accept-Language")
	languageCode := "en" // default
	if acceptLanguage != "" {
		parts := strings.Split(acceptLanguage, ",")
		if len(parts) > 0 {
			lang := strings.Split(parts[0], ";")[0]
			languageCode = strings.TrimSpace(lang)
		}
	}

	for _, round := range resp.Round {

		playName, err := s.Repo.GetPlayNameByPlayID(round.GameId, languageCode)

		if err != nil {
			log.Println("Error getting game name:", err)
			continue
		}

		log.Printf(
			"Game name for Game ID %d: %s\n",
			round.GameId,
			playName,
		)

		round.GameName = playName

		switch round.Status {
		case "start", "wait_bets", "accept_bets":
			round.Status = "Start Betting"

		case "roll_dice":
			round.Status = "Opening"

		case "settlement", "completed", "modified":
			round.Status = "Settled"

		}

	}
	log.Printf("Final gRPC Response with game names: %+v\n", resp)
	// result := map[string]any {
	// 	"response":    game[0]['result'][resp.Round.Id],
	// 	"page_no":    pagination.PageNo,
	// 	"page_size":  pagination.PageSize,
	// 	"total":      total,
	// 	"total_page": (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
	// }
	// bannerResources := make([]Resources.BannerResource, 0, total)
	// for i := range banners {
	// 	bannerResources = append(bannerResources, mapToBannerResource(&banners[i]))
	// }
	// config-service  | 2026/05/12 12:38:29 gRPC Response: round:{id:1  table_id:1  shoe_no:1  round_no:1  status:"completed"  betting_start_at:"2026-05-09 11:10:20"  betting_end_at:"2026-05-09 11:10:20"}  found:true

	total := resp.Total

	totalPage := (int64(total) + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	// 4. Return formatted paginated response
	// Map proto messages into plain maps so JSON marshaller includes empty strings
	records := make([]map[string]interface{}, 0, len(resp.Round))
	for _, r := range resp.Round {
		rec := map[string]interface{}{
			"id":               r.Id,
			"table_id":         r.TableId,
			"shoe_no":          r.ShoeNo,
			"round_no":         r.RoundNo,
			"status":           r.Status,
			"dice1":            r.Dice1,
			"dice2":            r.Dice2,
			"dice3":            r.Dice3,
			"betting_start_at": r.BettingStartAt,
			// ensure settled_at key is always present even when empty string
			"settled_at": r.SettledAt,
			"table_no":   r.TableNo,
			"game_id":    r.GameId,
			"game_name":  r.GameName,
			"remark":     r.Remark,
		}
		records = append(records, rec)
	}

	return map[string]interface{}{
		"records":    records,
		"page_no":    pagination.PageNo,
		"page_size":  pagination.PageSize,
		"total":      total,
		"total_page": totalPage,
	}, nil
}

func (s *ReportService) UpdateRoundRemark(c *fiber.Ctx, roundID uint64, remark string) (map[string]interface{}, error) {
	conn := config.GameEngineGrpcClient()
	defer conn.Close()

	client := commonpb.NewRoundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// fetch current round to capture original remark and other meta
	roundResp, err := client.GetRoundById(ctx, &commonpb.GetRoundByIdRequest{Id: roundID})
	log.Println("GetRoundById response:", roundResp)
	if err != nil {
		log.Println("GetRoundById error:", err)
		return nil, err
	}
	if roundResp == nil {
		return nil, fmt.Errorf("GetRoundById returned nil response for round id %d", roundID)
	}
	if !roundResp.Found || roundResp.Round == nil {
		return nil, fmt.Errorf("round not found for id %d", roundID)
	}

	originalRemark := roundResp.Round.Remark
	tableNo := roundResp.Round.TableNo
	shoeNo := roundResp.Round.ShoeNo
	roundNo := roundResp.Round.RoundNo
	playID := roundResp.Round.GameId

	resp, err := client.UpdateRoundRemark(ctx, &commonpb.UpdateRoundRemarkRequest{
		Id:     roundID,
		Remark: remark,
	})

	// prepare log entry
	now := time.Now()
	opResult := "Failure"
	if err == nil && resp != nil && resp.Success {
		opResult = "Success"
	}

	// try to get operator info from context if available
	operator := ""
	if u := c.Locals("user"); u != nil {
		if claims, ok := u.(jwt.MapClaims); ok && claims != nil {
			// prefer username/email/name, fall back to id or sub
			if v, exists := claims["email"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				}
			} else if v, exists := claims["name"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				}
			} else if v, exists := claims["sub"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				case float64:
					operator = fmt.Sprintf("%v", t)
				}
			} else if v, exists := claims["id"]; exists {
				switch t := v.(type) {
				case float64:
					operator = fmt.Sprintf("%d", int64(t))
				case string:
					operator = t
				}
			}
		}
		if operator == "" {
			operator = "unknown"
		}
	}

	rlog := &models.RoundRemarkLog{
		RoundID:         roundID,
		TableNo:         tableNo,
		ShoeNo:          shoeNo,
		RoundNo:         roundNo,
		GameID:          playID,
		OriginalRemark:  originalRemark,
		NewRemark:       remark,
		OperatorAccount: operator,
		OperationTime:   now,
		OperationIP:     c.IP(),
		OperationDevice: c.Get("User-Agent"),
		OperationResult: opResult,
	}

	if err2 := s.Repo.CreateRoundRemarkLog(rlog); err2 != nil {
		// log but do not fail the operation because of logging failure
		log.Println("failed to insert remark log:", err2)
	}

	if err != nil {
		return nil, err
	}
	acceptLanguage := c.Get("Accept-Language")
	languageCode := "en" // default
	if acceptLanguage != "" {
		parts := strings.Split(acceptLanguage, ",")
		if len(parts) > 0 {
			lang := strings.Split(parts[0], ";")[0]
			languageCode = strings.TrimSpace(lang)
		}
	}
	playName, err := s.Repo.GetPlayNameByPlayID(playID, languageCode)
	roundMap := map[string]interface{}{
		"tableNo":        tableNo,
		"shoeNo":         shoeNo,
		"roundNo":        roundNo,
		"gameName":       playName,
		"originalRemark": originalRemark,
	}

	return map[string]interface{}{
		"success":   resp.Success,
		"message":   resp.Message,
		"roundResp": roundMap,
	}, nil
}

func (s *ReportService) UpdateRoundResult(c *fiber.Ctx, roundID uint64, req *requests.UpdateRoundResultRequest) (map[string]interface{}, error) {
	conn := config.GameEngineGrpcClient()
	defer conn.Close()

	client := commonpb.NewRoundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// fetch current round to capture original result and other meta
	roundResp, err := client.GetRoundById(ctx, &commonpb.GetRoundByIdRequest{Id: roundID})
	if err != nil {
		log.Println("GetRoundById error:", err)
		return nil, err
	}
	if roundResp == nil {
		return nil, fmt.Errorf("GetRoundById returned nil response for round id %d", roundID)
	}
	if !roundResp.Found || roundResp.Round == nil {
		return nil, fmt.Errorf("round not found for id %d", roundID)
	}

	initialResult := formatDiceResult(roundResp.Round.OriginalDice1, roundResp.Round.OriginalDice2, roundResp.Round.OriginalDice3)
	currentResultBefore := formatDiceResult(roundResp.Round.Dice1, roundResp.Round.Dice2, roundResp.Round.Dice3)
	tableNo := roundResp.Round.TableNo
	shoeNo := roundResp.Round.ShoeNo
	roundNo := roundResp.Round.RoundNo
	playID := roundResp.Round.GameId

	resp, err := client.UpdateRoundResult(ctx, &commonpb.UpdateRoundResultRequest{
		Id:     roundID,
		Dice1:  req.Dice1,
		Dice2:  req.Dice2,
		Dice3:  req.Dice3,
		Remark: req.Remark,
		Status: "",
	})

	now := time.Now()
	opResult := "Failure"
	if err == nil && resp != nil && resp.Success {
		opResult = "Success"
	}

	operator := ""
	if u := c.Locals("user"); u != nil {
		if claims, ok := u.(jwt.MapClaims); ok && claims != nil {
			if v, exists := claims["email"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				}
			} else if v, exists := claims["name"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				}
			} else if v, exists := claims["sub"]; exists {
				switch t := v.(type) {
				case string:
					operator = t
				case float64:
					operator = fmt.Sprintf("%v", t)
				}
			} else if v, exists := claims["id"]; exists {
				switch t := v.(type) {
				case float64:
					operator = fmt.Sprintf("%d", int64(t))
				case string:
					operator = t
				}
			}
		}
		if operator == "" {
			operator = "unknown"
		}
	}

	modifiedResult := formatDiceResult(req.Dice1, req.Dice2, req.Dice3)

	rlog := &models.RoundResultModificationLog{
		RoundID:                   roundID,
		TableNo:                   tableNo,
		ShoeNo:                    shoeNo,
		RoundNo:                   roundNo,
		GameID:                    playID,
		InitialOriginalResult:     initialResult,
		CurrentResultBeforeModify: currentResultBefore,
		ModifiedResult:            modifiedResult,
		ModificationReason:        req.Remark,
		OperatorAccount:           operator,
		OperationTime:             now,
		OperationIP:               c.IP(),
		OperationDevice:           c.Get("User-Agent"),
		OperationResult:           opResult,
	}

	if err2 := s.Repo.CreateRoundResultModificationLog(rlog); err2 != nil {
		log.Println("failed to insert result modification log:", err2)
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	}, nil
}

func formatDiceResult(d1, d2, d3 uint32) string {
	return fmt.Sprintf("%d,%d,%d", d1, d2, d3)
}
func (s *ReportService) AuditAction(
	c *fiber.Ctx,
	action string,
	model string,
	user jwt.MapClaims,
	oldData any,
	newData any,
) {

	helper.PublishAudit(
		s.Audit,
		c,
		action,
		model,
		user,
		oldData,
		newData,
	)
}
