package Services

import (
	enums "config-service/Enums"
	"config-service/models"
	"config-service/pb/commonpb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AdminService struct {
	commonpb.UnimplementedAdminServiceServer
	DB *gorm.DB
}

type tableCatalogRow struct {
	GameID     uint64 `gorm:"column:game_id"`
	GameType   uint   `gorm:"column:game_data_id"`
	GameName   string `gorm:"column:game_name"`
	HallID     uint64 `gorm:"column:hall_id"`
	HallName   string `gorm:"column:hall_name"`
	HallSort   int32  `gorm:"column:hall_sort"`
	TableID    uint64 `gorm:"column:table_id"`
	TableNo    string `gorm:"column:table_no"`
	TableName  string `gorm:"column:table_name"`
	PlayType   uint32 `gorm:"column:play_type"`
	Status     string `gorm:"column:status"`
	AppIcon    string `gorm:"column:app_icon"`
	PcIcon     string `gorm:"column:pc_icon"`
	BetSecond  uint32 `gorm:"column:bet_second"`
	CalcSecond uint32 `gorm:"column:calc_second"`
	Mode       string `gorm:"column:mode"`

	Resolution int32
	RoomName   string
}
type playTableResolutionRow struct {
	Resolution int    `gorm:"column:resolution"`
	RoomName   string `gorm:"column:room_name"`
}

type redLimitPayload struct {
	PlayLimitDetails []*redLimitPlayDetail `json:"gameLimitDetails"`
}

type redLimitPlayDetail struct {
	GameName        string            `json:"gameName"`
	AutoExist       uint              `json:"gameAutoExits"`
	GameID          uint64            `json:"gameId"`
	Currency        string            `json:"currency"`
	LimitPlanVoList []*redLimitPlanVo `json:"limitPlanVoList"`
	Odds            []*redLimitOddsVo `json:"odds"`
}

type redLimitPlanVo struct {
	Name         string              `json:"name"`
	LimitID      uint64              `json:"limitId"`
	Status       uint32              `json:"status"`
	Sort         uint32              `json:"sort"`
	PlayType     uint32              `json:"playType"`
	DetailVoList []*redLimitDetailVo `json:"detailVoList"`
}

type redLimitDetailVo struct {
	PlayID   string  `json:"playId"`
	PlayName string  `json:"playName"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
}

type redLimitOddsVo struct {
	Tiger        float64 `json:"tiger"`
	Klouk        float64 `json:"klouk"`
	Fish         float64 `json:"fish"`
	Chicken      float64 `json:"chicken"`
	Shrimp       float64 `json:"shrimp"`
	Crab         float64 `json:"crab"`
	TripleSymbol float64 `json:"triple_symbol"`
	TripleColor  float64 `json:"triple_color"`
	DoubleAnimal float64 `json:"double_animal"`
	PlayType     uint32  `json:"playType"`
	GoldCoin     float64 `json:"gold_coin"`
}

type redLimitGroupBucket struct {
	Group          *redLimitPlayDetail
	OddsByPlayType map[uint32]*redLimitOddsVo
}

var redLimitPlayIDMap = map[string]string{
	"tiger":         "1",
	"klouk":         "2",
	"fish":          "3",
	"chicken":       "4",
	"shrimp":        "5",
	"crab":          "6",
	"triple_symbol": "7",
	"triple_color":  "8",
	"double_animal": "9",
	"gold_coin":     "10",
}

// Display order only — does NOT affect PlayID values
var redLimitPlayOrder = map[string]int{
	"tiger":         1,
	"fish":          2,
	"shrimp":        3, // Gourd
	"crab":          4,
	"chicken":       5,
	"gold_coin":     6,
	"klouk":         7,
	"double_animal": 8,
	"triple_symbol": 9,  // Double Dices
	"triple_color":  10, // keep last (not in your list)
}

type playTranslationRow struct {
	TableID      uint64 `gorm:"column:table_id"`
	LanguageCode string `gorm:"column:language_code"`
	Content      string `gorm:"column:content"`
}

type playTypeTranslationRow struct {
	MasterGameTypeID uint64 `gorm:"column:master_game_type_id"`
	LanguageCode     string `gorm:"column:language_code"`
	Name             string `gorm:"column:name"`
}

type tableResolutionRow struct {
	GameTableID uint64 `gorm:"column:game_table_id"`
	Resolution  int32  `gorm:"column:resolution"`
	RoomName    string `gorm:"column:room_name"`
}

func (s *AdminService) loadPlayNames(ctx context.Context, playIDs []uint64, lang string) (map[uint64]string, error) {
	out := make(map[uint64]string)
	if len(playIDs) == 0 {
		return out, nil
	}

	langs := preferredLanguages(lang)

	var rows []playTranslationRow
	if err := s.DB.WithContext(ctx).
		Table("game_translations").
		Select("table_id, language_code, content").
		Where("table_type = ?", "games").
		Where("lang_name = ?", "game_name").
		Where("table_id IN ?", playIDs).
		Where("language_code IN ?", langs).
		Find(&rows).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tmp := make(map[uint64]map[string]string)
	for _, row := range rows {
		if _, ok := tmp[row.TableID]; !ok {
			tmp[row.TableID] = make(map[string]string)
		}
		tmp[row.TableID][row.LanguageCode] = row.Content
	}

	for id, langMap := range tmp {
		out[id] = pickPreferredValue(langMap, lang)
	}

	return out, nil
}

func (s *AdminService) loadRedLimitNames(ctx context.Context, limitIDs []uint64, lang string) (map[uint64]string, error) {
	out := make(map[uint64]string)
	if len(limitIDs) == 0 {
		return out, nil
	}

	langs := preferredLanguages(lang)

	var rows []playTranslationRow
	if err := s.DB.WithContext(ctx).
		Table("game_translations").
		Select("table_id, language_code, content").
		Where("table_type = ?", "game_red_limits").
		Where("lang_name = ?", "red_limit").
		Where("table_id IN ?", limitIDs).
		Where("language_code IN ?", langs).
		Find(&rows).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tmp := make(map[uint64]map[string]string)
	for _, row := range rows {
		if _, ok := tmp[row.TableID]; !ok {
			tmp[row.TableID] = make(map[string]string)
		}
		tmp[row.TableID][row.LanguageCode] = row.Content
	}

	for id, langMap := range tmp {
		out[id] = pickPreferredValue(langMap, lang)
	}

	return out, nil
}

func (s *AdminService) loadMasterPlayTypeNames(ctx context.Context, masterIDs []uint64, lang string) (map[uint64]string, error) {
	out := make(map[uint64]string)
	if len(masterIDs) == 0 {
		return out, nil
	}

	langs := preferredLanguages(lang)

	var rows []playTypeTranslationRow
	if err := s.DB.WithContext(ctx).
		Table("game_type_translations").
		Select("master_game_type_id, language_code, name").
		Where("master_game_type_id IN ?", masterIDs).
		Where("language_code IN ?", langs).
		Find(&rows).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tmp := make(map[uint64]map[string]string)
	for _, row := range rows {
		if _, ok := tmp[row.MasterGameTypeID]; !ok {
			tmp[row.MasterGameTypeID] = make(map[string]string)
		}
		tmp[row.MasterGameTypeID][row.LanguageCode] = row.Name
	}

	for id, langMap := range tmp {
		out[id] = pickPreferredValue(langMap, lang)
	}

	return out, nil
}

func preferredLanguages(lang string) []string {
	if lang == "" || lang == "en" {
		return []string{"en"}
	}
	return []string{lang, "en"}
}

func pickPreferredValue(values map[string]string, lang string) string {
	if values == nil {
		return ""
	}
	if v := strings.TrimSpace(values[lang]); v != "" {
		return v
	}
	if v := strings.TrimSpace(values["en"]); v != "" {
		return v
	}
	return ""
}

func mapKeysUint64(m map[uint64]struct{}) []uint64 {
	out := make([]uint64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
func (s *AdminService) loadMasterPlayTypesByPlayType(ctx context.Context, playTypeIDs []uint64) (map[uint64][]models.MasterGameType, error) {
	result := make(map[uint64][]models.MasterGameType)
	if len(playTypeIDs) == 0 {
		return result, nil
	}

	var masters []models.MasterGameType
	if err := s.DB.WithContext(ctx).
		Model(&models.MasterGameType{}).
		Where("game_type_id IN ?", playTypeIDs).
		Find(&masters).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, master := range masters {
		result[master.GameTypeID] = append(result[master.GameTypeID], master)
	}

	for playTypeID := range result {
		sortMasterPlayTypes(result[playTypeID])
	}

	return result, nil
}

func (s *AdminService) loadPlayTypeDataByPlayType(
	ctx context.Context,
	playTypeIDs []uint64,
	playIDs []uint64,
) (map[string][]models.GameTypeData, error) {

	result := make(map[string][]models.GameTypeData)

	if len(playTypeIDs) == 0 || len(playIDs) == 0 {
		return result, nil
	}

	var rows []models.GameTypeData

	err := s.DB.WithContext(ctx).
		Model(&models.GameTypeData{}).
		Preload("MasterGameType").
		Where("game_id IN ?", playIDs).
		Where("master_game_type_id IN (?)",
			s.DB.
				Model(&models.MasterGameType{}).
				Select("id").
				Where("game_type_id IN ?", playTypeIDs),
		).
		Find(&rows).Error

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, row := range rows {

		if row.MasterGameType == nil {
			continue
		}

		playTypeID := row.MasterGameType.GameTypeID

		key := fmt.Sprintf("%d_%d", row.GameID, playTypeID)

		result[key] = append(result[key], row)
	}

	for key := range result {
		sortPlayTypeData(result[key])
	}

	return result, nil
}

// func sortGameTypeData(items []models.GameTypeData) {

// 	orderMap := map[string]int{
// 		"tiger":         1,
// 		"klouk":         2,
// 		"fish":          3,
// 		"chicken":       4,
// 		"shrimp":        5,
// 		"crab":          6,
// 		"triple_symbol": 7,
// 		"triple_color":  8,
// 		"double_animal": 9,
// 		"gold_coin":     10,
// 	}

// 	for i := 0; i < len(items); i++ {
// 		for j := i + 1; j < len(items); j++ {

// 			left := 999
// 			right := 999

// 			if items[i].MasterGameType != nil {
// 				left = orderMap[strings.TrimSpace(
// 					strings.ToLower(items[i].MasterGameType.KeyName),
// 				)]
// 			}

// 			if items[j].MasterGameType != nil {
// 				right = orderMap[strings.TrimSpace(
// 					strings.ToLower(items[j].MasterGameType.KeyName),
// 				)]
// 			}

// 			if left == 0 {
// 				left = 999
// 			}

// 			if right == 0 {
// 				right = 999
// 			}

//				if left > right {
//					items[i], items[j] = items[j], items[i]
//				}
//			}
//		}
//	}
func sortPlayTypeData(items []models.GameTypeData) {

	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {

			left := 999
			right := 999

			if items[i].MasterGameType != nil {
				left = redLimitPlayOrder[strings.TrimSpace(
					strings.ToLower(items[i].MasterGameType.KeyName),
				)]
			}

			if items[j].MasterGameType != nil {
				right = redLimitPlayOrder[strings.TrimSpace(
					strings.ToLower(items[j].MasterGameType.KeyName),
				)]
			}

			if left == 0 {
				left = 999
			}

			if right == 0 {
				right = 999
			}

			if left > right {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func (s *AdminService) buildRedLimitPayload(ctx context.Context, playID uint64, language string, currency string) (*redLimitPayload, error) {
	lang := normalizeLanguage(language)
	currencyCode := enums.ParseCurrencyID(currency)
	var limits []models.GameRedLimit
	query := s.DB.WithContext(ctx).
		Model(&models.GameRedLimit{}).
		Select("game_red_limits.*, games.no_auto_exits").
		Joins("JOIN games ON games.id = game_red_limits.game_type").
		Where("game_red_limits.status = ?", "active").
		Where("games.status = ?", "enable").
		Where("game_red_limits.game_type_currency = ?", currencyCode).
		Preload("RedLimitTable").
		Order("game_red_limits.game_type ASC").
		Order("game_red_limits.game_type_currency ASC").
		Order("game_red_limits.sort ASC").
		Order("game_red_limits.id ASC")

	if playID > 0 {
		query = query.Where("game_red_limits.game_type = ?", playID)
	}

	if err := query.Find(&limits).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payload := &redLimitPayload{
		PlayLimitDetails: make([]*redLimitPlayDetail, 0),
	}

	if len(limits) == 0 {
		return payload, nil
	}

	playIDSet := make(map[uint64]struct{})
	limitIDSet := make(map[uint64]struct{})
	playTypeIDSet := make(map[uint64]struct{})

	for _, limit := range limits {
		playIDSet[uint64(limit.GameType)] = struct{}{}
		limitIDSet[uint64(limit.ID)] = struct{}{}
		playTypeIDSet[uint64(limit.GamePlayType)] = struct{}{}
	}

	playIDs := mapKeysUint64(playIDSet)
	limitIDs := mapKeysUint64(limitIDSet)
	playTypeIDs := mapKeysUint64(playTypeIDSet)

	playTypeDataByPlayType, err := s.loadPlayTypeDataByPlayType(
		ctx,
		playTypeIDs,
		playIDs,
	)

	if err != nil {
		return nil, err
	}

	playNames, err := s.loadPlayNames(ctx, playIDs, lang)
	if err != nil {
		return nil, err
	}

	limitNames, err := s.loadRedLimitNames(ctx, limitIDs, lang)
	if err != nil {
		return nil, err
	}

	masterIDs := make([]uint64, 0)
	masterIDSeen := make(map[uint64]struct{})

	for _, items := range playTypeDataByPlayType {

		for _, item := range items {

			if item.MasterGameType == nil {
				continue
			}

			masterID := item.MasterGameType.ID

			if _, ok := masterIDSeen[masterID]; ok {
				continue
			}

			masterIDSeen[masterID] = struct{}{}
			masterIDs = append(masterIDs, masterID)
		}
	}

	masterPlayNames, err := s.loadMasterPlayTypeNames(ctx, masterIDs, lang)
	if err != nil {
		return nil, err
	}

	groupBuckets := make(map[string]*redLimitGroupBucket)

	for _, limit := range limits {
		groupKey := strconv.FormatUint(uint64(limit.GameType), 10) + "_" + strconv.FormatUint(uint64(limit.GameTypeCurrency), 10)

		bucket, exists := groupBuckets[groupKey]
		if !exists {
			playName := playNames[uint64(limit.GameType)]
			if playName == "" {
				playName = "Game"
			}

			group := &redLimitPlayDetail{
				GameName:        playName,
				AutoExist:       limit.GameAutoExits,
				GameID:          uint64(limit.GameType),
				Currency:        s.resolveCurrencyCode(uint64(limit.GameTypeCurrency)),
				LimitPlanVoList: make([]*redLimitPlanVo, 0),
				Odds:            make([]*redLimitOddsVo, 0),
			}

			bucket = &redLimitGroupBucket{
				Group:          group,
				OddsByPlayType: make(map[uint32]*redLimitOddsVo),
			}

			groupBuckets[groupKey] = bucket
			payload.PlayLimitDetails = append(payload.PlayLimitDetails, group)
		}

		planName := limitNames[uint64(limit.ID)]
		if planName == "" {
			planName = "Standard"
		}

		plan := &redLimitPlanVo{
			Name:         planName,
			LimitID:      uint64(limit.ID),
			Status:       redLimitStatusToInt(limit.Status),
			Sort:         uint32(limit.Sort),
			PlayType:     uint32(limit.GamePlayType),
			DetailVoList: make([]*redLimitDetailVo, 0),
		}

		detailMap := make(map[string]models.RedLimitTable, len(limit.RedLimitTable))
		for _, d := range limit.RedLimitTable {
			key := strings.TrimSpace(strings.ToLower(d.GameRedLimitKey))
			detailMap[key] = d
		}

		oddsVo, exists := bucket.OddsByPlayType[uint32(limit.GamePlayType)]
		if !exists {
			oddsVo = &redLimitOddsVo{
				PlayType: uint32(limit.GamePlayType),
			}
			bucket.OddsByPlayType[uint32(limit.GamePlayType)] = oddsVo
			bucket.Group.Odds = append(bucket.Group.Odds, oddsVo)
		}

		mapKey := fmt.Sprintf(
			"%d_%d",
			uint64(limit.GameType),
			uint64(limit.GamePlayType),
		)

		items := playTypeDataByPlayType[mapKey]

		for _, item := range items {

			if item.MasterGameType == nil {
				continue
			}

			master := item.MasterGameType

			key := strings.TrimSpace(strings.ToLower(master.KeyName))

			playName := masterPlayNames[master.ID]

			if playName == "" {
				playName = humanizePlayKey(key)
			}

			detailRow, found := detailMap[key]

			minVal := 0.0
			maxVal := 0.0

			if found {
				minVal = float64(detailRow.MinLimit)
				maxVal = float64(detailRow.MaxLimit)
			}

			plan.DetailVoList = append(plan.DetailVoList, &redLimitDetailVo{
				PlayID:   redLimitPlayIDMap[key],
				PlayName: key,
				Min:      minVal,
				Max:      maxVal,
			})

			// odds now come from game_type_data table
			applyRedLimitOdds(oddsVo, key, float64(item.Odds))
		}

		bucket.Group.LimitPlanVoList = append(bucket.Group.LimitPlanVoList, plan)
	}

	return payload, nil
}

func (s *AdminService) resolveCurrencyCode(currencyID uint64) string {

	return enums.ParseCurrency(currencyID)

}

func redLimitStatusToInt(status string) uint32 {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active", "enable", "enabled":
		return 1
	default:
		return 0
	}
}

func applyRedLimitOdds(target *redLimitOddsVo, key string, value float64) {
	switch key {
	case "tiger":
		target.Tiger = value
	case "klouk":
		target.Klouk = value
	case "fish":
		target.Fish = value
	case "chicken":
		target.Chicken = value
	case "shrimp":
		target.Shrimp = value
	case "crab":
		target.Crab = value
	case "triple_symbol":
		target.TripleSymbol = value
	case "triple_color":
		target.TripleColor = value
	case "double_animal":
		target.DoubleAnimal = value
	case "gold_coin":
		target.GoldCoin = value
	}
}

func humanizePlayKey(key string) string {
	switch key {
	case "triple_symbol":
		return "Any Triple Symbols"
	case "triple_color":
		return "Any Triple Colors"
	case "double_animal":
		return "Double Animal"
	default:
		parts := strings.Split(key, "_")
		for i, part := range parts {
			if part == "" {
				continue
			}
			if len(part) == 1 {
				parts[i] = strings.ToUpper(part)
			} else {
				parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
		}
		return strings.Join(parts, " ")
	}
}

// func sortMasterGameTypes(items []models.MasterGameType) {
// 	orderMap := map[string]int{
// 		"tiger":         1,
// 		"klouk":         2,
// 		"fish":          3,
// 		"chicken":       4,
// 		"shrimp":        5,
// 		"crab":          6,
// 		"triple_symbol": 7,
// 		"triple_color":  8,
// 		"double_animal": 9,
// 		"gold_coin":     10,
// 	}

// 	for i := 0; i < len(items); i++ {
// 		for j := i + 1; j < len(items); j++ {
// 			left := orderMap[strings.TrimSpace(strings.ToLower(items[i].KeyName))]
// 			right := orderMap[strings.TrimSpace(strings.ToLower(items[j].KeyName))]

// 			if left == 0 {
// 				left = 999
// 			}
// 			if right == 0 {
// 				right = 999
// 			}

//				if left > right {
//					items[i], items[j] = items[j], items[i]
//				}
//			}
//		}
//	}
func sortMasterPlayTypes(items []models.MasterGameType) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			left := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[i].KeyName))]
			right := redLimitPlayOrder[strings.TrimSpace(strings.ToLower(items[j].KeyName))]

			if left == 0 {
				left = 999
			}
			if right == 0 {
				right = 999
			}

			if left > right {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// GetTables implements the streaming RPC you defined
func (s *AdminService) GetTables(req *commonpb.Empty, stream commonpb.AdminService_GetTablesServer) error {
	type TableWithPlayStatus struct {
		models.GameTable
		PlayStatus string `gorm:"column:game_status"`
		HallStatus string `gorm:"column:hall_status"`
	}

	var tables []TableWithPlayStatus

	err := s.DB.
		Table("game_tables AS gt").
		Select(`
			gt.*,
			g.status AS game_status,
			gh.status AS hall_status
		`).
		Joins("LEFT JOIN games AS g ON g.id = gt.game_type").
		Joins("LEFT JOIN game_halls AS gh ON gh.id = gt.affiliate_office").
		Scan(&tables).Error

	if err != nil {
		return err
	}

	for _, t := range tables {
		tableStatus := t.Status

		//If parent game is disabled, force table disabled
		if t.PlayStatus != "enable" || t.HallStatus != "enable" {
			tableStatus = "disable"
		}

		resp := &commonpb.GetTablesResponse{
			Id:             uint64(t.ID),
			TableNo:        t.TableNumber,
			HallId:         uint64(t.AffiliateOffice),
			GameId:         uint64(t.GameType),
			Mode:           t.Mode,
			BettingTime:    uint32(t.BettingTime),
			SettlementTime: uint32(t.SettlementTime),
			Status:         tableStatus,
			KafkaTopic:     t.Topic,
			GameTypeId:     uint64(t.GameTypeID),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

func normalizeLanguage(lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))

	switch lang {
	case "en":
		return "en"
	case "ch", "zh", "zh-cn", "cn":
		return "ch"
	case "kh", "km":
		return "kh"
	default:
		return "en"
	}
}

func (s *AdminService) GetEngineStatus(ctx context.Context, req *commonpb.Empty) (*commonpb.GameEngineStatus, error) {
	disablePlayEngine, err := strconv.ParseBool(os.Getenv("DISABLE_GAME_ENGINE"))

	log.Println("game engine status", disablePlayEngine)
	if err != nil {
		disablePlayEngine = false
	}

	resp := &commonpb.GameEngineStatus{
		Status: disablePlayEngine,
	}

	return resp, nil
}

func (s *AdminService) GetTableCatalog(ctx context.Context, req *commonpb.TableCatalogRequest) (*commonpb.TableCatalogResponse, error) {
	log.Println("GetTableCatalog admin service")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	lang := normalizeLanguage(req.GetLanguage())

	var rows []tableCatalogRow

	query := s.DB.WithContext(ctx).
		Table("game_tables as gt").
		Select(`
			COALESCE(g.id, 0) as game_id,
			COALESCE(g.game_data_id, 0) as game_data_id,
			COALESCE(gtr.content, '') as game_name,
			COALESCE(gh.id, 0) as hall_id,
			COALESCE(htr.content, '') as hall_name,
			COALESCE(gh.sort_order, 0) as hall_sort,
			COALESCE(gt.id, 0) as table_id,
			COALESCE(gt.table_number, '') as table_no,
			COALESCE(ttr.content, '') as table_name,
			COALESCE(gt.game_type_id, 0) as play_type,
			COALESCE(gt.status, '') as status,
			COALESCE(gt.mobile_icon, '') as app_icon,
			COALESCE(gt.pc_icon, '') as pc_icon,
			COALESCE(gt.betting_time, 0) as bet_second,
			COALESCE(gt.settlement_time, 0) as calc_second,
			COALESCE(gt.mode, '') as mode
		`).
		Joins("LEFT JOIN games g ON g.id = gt.game_type").
		Joins("LEFT JOIN game_halls gh ON gh.id = gt.affiliate_office").

		// game translation
		Joins(`
			LEFT JOIN game_translations as gtr
				ON gtr.table_id = g.id
				AND gtr.table_type = ?
				AND gtr.lang_name = ?
				AND gtr.language_code = ?
		`, "games", "game_name", lang).

		// hall translation
		Joins(`
			LEFT JOIN game_translations as htr
				ON htr.table_id = gh.id
				AND htr.table_type = ?
				AND htr.lang_name = ?
				AND htr.language_code = ?
		`, "game_halls", "game_hall_name", lang).

		// table translation
		Joins(`
			LEFT JOIN game_translations as ttr
				ON ttr.table_id = gt.id
				AND ttr.table_type = ?
				AND ttr.lang_name = ?
				AND ttr.language_code = ?
		`, "game_tables", "game_table_name", lang)

	if req.GetGameId() > 0 {
		query = query.Where("gt.game_type = ?", req.GetGameId()) //this may change in future
	}

	if req.GetTableId() != "" {
		query = query.Where("gt.table_number = ?", req.GetTableId())
	}

	err := query.
		Order("g.id ASC").
		Order("gh.sort_order ASC").
		Order("gt.table_number ASC").
		Where("gt.status = ?", "active").
		Where("g.status = ?", "enable").
		Where("gh.status = ?", "enable").
		Scan(&rows).Error
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Println("GetTableCatalog rows", rows)

	tableIDs := make([]uint64, 0, len(rows))

	for _, row := range rows {
		tableIDs = append(tableIDs, row.TableID)
	}

	resp := &commonpb.TableCatalogResponse{
		Tables: make([]*commonpb.TableCatalogItem, 0, len(rows)),
	}
	log.Println("GetTableCatalog resp", resp)

	var resolutionRows []tableResolutionRow

	if len(tableIDs) > 0 {
		err = s.DB.WithContext(ctx).
			Table("game_table_resolutions").
			Select(`
				game_table_id,
				resolution,
				room_name
			`).
			Where("game_table_id IN ?", tableIDs).
			Order("id ASC").
			Scan(&resolutionRows).Error

		if err != nil {
			log.Println(
				"GetTableCatalog resolutions query error:",
				err,
			)

			return nil, status.Error(
				codes.Internal,
				err.Error(),
			)
		}
	}

	log.Println(
		"GetTableCatalog resolutionRows:",
		resolutionRows,
	)

	resolutionsByTable := make(
		map[uint64][]*commonpb.ResolutionList,
	)

	for _, resolutionRow := range resolutionRows {

		resolution := &commonpb.ResolutionList{
			ResolutionType: resolutionRow.Resolution,
			RoomName:       resolutionRow.RoomName,
		}

		resolutionsByTable[resolutionRow.GameTableID] =
			append(
				resolutionsByTable[resolutionRow.GameTableID],
				resolution,
			)
	}

	// ------------------------------------------------------------
	// 6. Build response
	// ------------------------------------------------------------

	for _, row := range rows {
		log.Println("rows list", row)
		var mobileIcon, pcIcon string
		if row.AppIcon != "" {
			mobileIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + row.AppIcon
		}

		if row.PcIcon != "" {
			pcIcon = os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + row.PcIcon
		}
		tableResolutions := resolutionsByTable[row.TableID]

		if tableResolutions == nil {
			tableResolutions = make(
				[]*commonpb.ResolutionList,
				0,
			)
		}

		//for resolution

		resp.Tables = append(resp.Tables, &commonpb.TableCatalogItem{
			GameId:         row.GameID,
			GameType:       uint64(row.GameType),
			GameName:       row.GameName,
			HallId:         row.HallID,
			HallName:       row.HallName,
			HallSort:       row.HallSort,
			TableId:        row.TableID,
			TableNo:        row.TableNo,
			TableName:      row.TableName,
			PlayType:       row.PlayType,
			Status:         row.Status,
			AppIcon:        mobileIcon,
			PcIcon:         pcIcon,
			BetSecond:      row.BetSecond,
			CalcSecond:     row.CalcSecond,
			Mode:           row.Mode,
			ResolutionList: tableResolutions,
		})
	}

	if req.GetTableId() != "" && len(rows) == 0 && len(resp.Tables) == 0 {
		return nil, status.Error(codes.NotFound, "table not found")

	}

	log.Println("table catalog response", resp)

	return resp, nil
}

func (s *AdminService) GetRedLimit(ctx context.Context, req *commonpb.RedLimitRequest) (*commonpb.RedLimitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	data, err := s.buildRedLimitPayload(ctx, req.GetGameId(), req.GetLanguage(), req.GetCurrency())
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &commonpb.RedLimitResponse{
		Datas: b,
	}, nil
}
func (s *AdminService) GetRedLimitByLimitID(ctx context.Context, req *commonpb.RedLimitRequestById) (*commonpb.RedLimitByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	log.Println("hello admin")
	if req.GetLimitId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit_id is required")
	}

	data, err := s.RedLimitForBetPayload(ctx, req.GetLimitId())
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &commonpb.RedLimitByIdResponse{
		Datas: b,
	}, nil
}

func (s *AdminService) RedLimitForBetPayload(ctx context.Context, id uint64) (any, error) {
	var limit models.GameRedLimit
	log.Println("s.DB", s.DB)
	err := s.DB.
		Preload("Translations").
		Preload("RedLimitTable").
		First(&limit, id).Error

	if err != nil {
		return nil, err
	}

	return limit, nil
}

// ===============================
// ADMIN SERVICE
// ===============================
func (s *AdminService) GetOddsBatch(
	ctx context.Context,
	req *commonpb.GetOddsBatchRequest,
) (*commonpb.GetOddsBatchResponse, error) {

	type OddsRow struct {
		BetCode          string
		MasterGameTypeID uint64
		Odds             float64
	}

	var rows []OddsRow

	// ====================================
	// COLLECT BET CODES
	// ====================================

	var betCodes []string

	for _, item := range req.Items {
		betCodes = append(
			betCodes,
			item.BetCode,
		)
	}

	// ====================================
	// GET GAME PLAY TYPE
	// ====================================

	var gamePlayType uint64

	if len(req.Items) > 0 {
		gamePlayType = req.Items[0].GamePlayType
	}

	// ====================================
	// SINGLE QUERY
	// ====================================

	err := s.DB.Debug().
		Table("master_game_types mgt").
		Select(`
			mgt.key_name as bet_code,
			mgt.id as master_game_type_id,
			gtd.odds as odds
		`).
		Joins(`
			INNER JOIN game_type_data gtd
			ON gtd.master_game_type_id = mgt.id
		`).
		Where(`
			mgt.key_name IN ?
			AND mgt.game_type_id = ?
		`,
			betCodes,
			gamePlayType,
		).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// ====================================
	// BUILD RESPONSE
	// ====================================

	var results []*commonpb.OddsResult

	for _, row := range rows {

		results = append(results, &commonpb.OddsResult{
			BetCode: row.BetCode,
			Odds:    float64(row.Odds),
		})
	}

	return &commonpb.GetOddsBatchResponse{
		Results: results,
	}, nil
}
func (s *AdminService) UpdateDealerTableOnlineStatus(
	ctx context.Context,
	req *commonpb.UpdateDealerTableOnlineStatusRequest,
) (*commonpb.UpdateDealerTableOnlineStatusResponse, error) {
	log.Println("UpdateDealerTableOnlineStatus called with request:", req)
	result := s.DB.
		Table("game_tables").
		Where("table_number = ?", req.TableNo).
		Where("game_type = ?", req.GameId).
		Update("kafka_online_status", req.TableOnlineStatus)

	if result.Error != nil {
		return &commonpb.UpdateDealerTableOnlineStatusResponse{
			Success: false,
			Message: result.Error.Error(),
		}, result.Error
	}

	if result.RowsAffected == 0 {
		return &commonpb.UpdateDealerTableOnlineStatusResponse{
			Success: false,
			Message: "table not found",
		}, nil
	}

	return &commonpb.UpdateDealerTableOnlineStatusResponse{
		Success: true,
		Message: "table online status updated successfully",
	}, nil
}
