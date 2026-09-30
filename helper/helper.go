package helper

import (
	"config-service/config"
	"config-service/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ParseStringToUint(str string) uint64 {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic("Failed to parse string to uint64")
	}
	return num
}

type CustomTime time.Time

// PaginationParams holds pagination info
// PaginationParams holds pagination info
type PaginationParams struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
	Offset   int `json:"-"`
	Limit    int `json:"-"`
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*ct = CustomTime(time.Time{})
		return nil
	}

	layout := "2006-01-02 15:04:05"
	s := strings.Trim(string(b), `"`)

	// Use Local timezone to prevent conversion to UTC
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return err
	}

	*ct = CustomTime(t)
	return nil
}

func (ct *CustomTime) ToTime() *time.Time {
	if ct == nil {
		return nil
	}
	t := time.Time(*ct)
	return &t
}

// GetPagination parses query params and returns PaginationParams
func GetPagination(c *fiber.Ctx) PaginationParams {
	pageNo, err := strconv.Atoi(c.Query("page_no", "1"))
	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (pageNo - 1) * pageSize

	return PaginationParams{
		PageNo:   pageNo,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

func NormalizeLanguage(raw string) string {
	lang := strings.ToLower(strings.TrimSpace(raw))
	lang = strings.Trim(lang, `"'`)
	lang = strings.ReplaceAll(lang, "_", "-")

	switch {
	case lang == "":
		return "en"
	case strings.HasPrefix(lang, "en"):
		return "en"
	case strings.HasPrefix(lang, "km"), strings.HasPrefix(lang, "kh"):
		return "kh"
	case strings.HasPrefix(lang, "cz"),
		strings.HasPrefix(lang, "ch"),
		strings.HasPrefix(lang, "zh"),
		strings.HasPrefix(lang, "cn"):
		return "ch"
	default:
		return "en"
	}
}

func GetTablesByIDs(tableIDs []uint, hallID uint, gameId uint) ([]models.GameTable, error) {
	var tables []models.GameTable
	sql := config.DB.
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
