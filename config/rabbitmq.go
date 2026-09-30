package config

import (
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn    *amqp.Connection
	RabbitChannel *amqp.Channel

	rabbitMutex sync.RWMutex
)

func ConnectRabbitMQ() {
	go reconnectLoop()
}

func reconnectLoop() {
	for {
		if err := connect(); err != nil {
			log.Printf("RabbitMQ connect failed: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("✅ Connected to RabbitMQ")

		notifyClose := RabbitConn.NotifyClose(make(chan *amqp.Error))

		err := <-notifyClose

		rabbitMutex.Lock()
		if RabbitChannel != nil {
			_ = RabbitChannel.Close()
			RabbitChannel = nil
		}

		if RabbitConn != nil {
			_ = RabbitConn.Close()
			RabbitConn = nil
		}
		rabbitMutex.Unlock()

		if err != nil {
			log.Printf("RabbitMQ disconnected: %v", err)
		}

		log.Println("Reconnecting RabbitMQ...")
		time.Sleep(5 * time.Second)
	}
}

func connect() error {

	url := "amqp://" +
		os.Getenv("RABBITMQ_USERNAME") + ":" +
		os.Getenv("RABBITMQ_PASSWORD") + "@" +
		os.Getenv("RABBITMQ_HOST") + ":" +
		os.Getenv("RABBITMQ_PORT") + "/"

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	_, err = ch.QueueDeclare(
		"audit_logs",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	rabbitMutex.Lock()
	RabbitConn = conn
	RabbitChannel = ch
	rabbitMutex.Unlock()

	return nil
}
func GetRabbitChannel() *amqp.Channel {
	rabbitMutex.RLock()
	defer rabbitMutex.RUnlock()
	return RabbitChannel
}
