package config

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQManager struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func InitRabbitMQManager() *RabbitMQManager {
	manager := &RabbitMQManager{}
	err := manager.Reconnect()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQManager: %v", err)
	}
	return manager
}

func (r *RabbitMQManager) Reconnect() error {
	// Close existing resources if open
	if r.Connection != nil {
		r.Connection.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}

	// Reconnect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Failed to reconnect to RabbitMQ:", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to reopen RabbitMQ channel:", err)
		conn.Close()
		return err
	}

	r.Connection = conn
	r.Channel = ch
	return nil
}
