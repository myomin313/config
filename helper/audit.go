// package helper

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type AuditService struct {
// 	Redis *redis.Client
// }
// // Generic publish function
// func (a *AuditService) Publish(eventType string, data any, ttl time.Duration) {
// 	bytes, err := json.Marshal(data)
// 	if err != nil {
// 		log.Println("audit marshal error:", err)
// 		return
// 	}

// 	key := fmt.Sprintf("audit:%s:%d", eventType, time.Now().UnixNano())
// 	log.Println("publish key", key)
// 	err = a.Redis.Set(context.Background(), key, bytes, ttl).Err()
// 	if err != nil {
// 		log.Printf("failed to publish audit log (%s): %v\n", eventType, err)
// 	}
// }

package helper

import (
	"context"
	"encoding/json"
	"log"

	"config-service/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuditService struct {
	Channel *amqp.Channel
	Queue   string
}

func (a *AuditService) Publish(eventType string, data any) {

	body, err := json.Marshal(data)
	if err != nil {
		log.Println("audit marshal error:", err)
		return
	}

	ch := config.GetRabbitChannel()

	if ch == nil {
		log.Println("RabbitMQ channel unavailable")
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",      // default exchange
		a.Queue, // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Type:         eventType,
			Body:         body,
		},
	)

	if err != nil {
		log.Printf("failed to publish audit log (%s): %v\n", eventType, err)
		return
	}

	log.Printf("✅ Published audit event: %s\n", eventType)
}
