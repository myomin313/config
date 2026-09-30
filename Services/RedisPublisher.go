package Services

import (
	"config-service/config"
	"config-service/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var tableDataChannel = "Kla-Klouk-game_table_data_" + os.Getenv("APP_ENV")

var GlobalVariablePublisher Publisher

type Event struct {
	// Version int       `json:"version"`
	EventID string `json:"event_id"`
	Event   any    `json:"event"`
	Channel string `json:"channel"`
	Source  string `json:"source"`
	Data    any    `json:"data,omitempty"`
	SentAt  int64  `json:"sent_at"`
}

func NewEvent(channel string, event any, data any) Event {
	return Event{
		// Version: 1,
		EventID: newEventID(),
		Event:   event,
		Channel: channel,
		Source:  "admin-service",
		Data:    data,
		SentAt:  time.Now().UTC().UnixMilli(),
	}
}

func newEventID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

type Publisher interface {
	Publish(ctx context.Context, evt Event) error
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Close() error
}

type RedisPublisher struct {
	client  *redis.Client
	timeout time.Duration
}

func NewRedisPublisher(client *redis.Client, timeout time.Duration) *RedisPublisher {
	if timeout <= 0 {
		timeout = 300 * time.Millisecond
	}

	return &RedisPublisher{
		client:  client,
		timeout: timeout,
	}
}
func (p *RedisPublisher) Publish(ctx context.Context, evt Event) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, evt.Channel, payload).Err()
}

func (p *RedisPublisher) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.client.Set(ctx, key, value, ttl).Err()
}

func (p *RedisPublisher) Close() error {
	return p.client.Close()
}

func (p *RedisPublisher) Get(
	ctx context.Context,
	key string,
) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	result, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func emit(publisher Publisher, evt Event) {
	if publisher == nil {
		return
	}

	if err := publisher.Publish(context.Background(), evt); err != nil {
		config.Error(fmt.Sprintf(
			"redis publish failed channel=%s event=%s err=%v",
			evt.Channel, evt.Event, err,
		))
	}
}

func (p *RedisPublisher) UpdateTableInfo(tbl *models.GameTable) {
	tableStatus := "active"

	if tbl.Status != "active" {
		tableStatus = "disable"
	}

	data := map[string]any{
		"id":              tbl.ID,
		"table_no":        tbl.TableNumber,
		"hall_id":         tbl.AffiliateOffice,
		"game_id":         tbl.GameType,
		"mode":            tbl.Mode, //"digital",
		"betting_time":    uint32(tbl.BettingTime),
		"settlement_time": uint32(tbl.SettlementTime),
		"status":          tableStatus,
		"kafka_topic":     tbl.Topic,
		"game_type_id":    tbl.GameTypeID,
	}

	emit(p, NewEvent(tableDataChannel, "table_update", data))
}

func (p *RedisPublisher) DisableTableInfo(tbl *models.GameTable) {

	data := map[string]any{
		"id":              tbl.ID,
		"table_no":        tbl.TableNumber,
		"hall_id":         tbl.AffiliateOffice,
		"game_id":         tbl.GameType,
		"mode":            tbl.Mode, // "digital",
		"betting_time":    uint32(tbl.BettingTime),
		"settlement_time": uint32(tbl.SettlementTime),
		"status":          "disable",
		"kafka_topic":     tbl.Topic,
		"game_type_id":    tbl.GameTypeID,
	}

	emit(p, NewEvent(tableDataChannel, "table_update", data))
}

func BulkUpdateTableInfo(tables []models.GameTable, status string) {

	for _, tbl := range tables {
		tableStatus := tbl.Status

		if status != "enable" && status != "active" {
			tableStatus = "disable"
		}

		data := map[string]any{
			"id":              tbl.ID,
			"table_no":        tbl.TableNumber,
			"hall_id":         tbl.AffiliateOffice,
			"game_id":         tbl.GameType,
			"mode":            tbl.Mode, //"digital",
			"betting_time":    uint32(tbl.BettingTime),
			"settlement_time": uint32(tbl.SettlementTime),
			"status":          tableStatus,
			"kafka_topic":     tbl.Topic,
			"game_type_id":    tbl.GameTypeID,
		}
		emit(GlobalVariablePublisher, NewEvent(tableDataChannel, "table_update", data))
	}

}
