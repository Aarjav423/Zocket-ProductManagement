package queue

import (
	"encoding/json"
	"log"
	"time"

	"project/config"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func StartImageProcessingConsumer(db *gorm.DB, manager *config.RabbitMQManager) {
	var ch *amqp.Channel
	var closeChan <-chan *amqp.Error

	for {
		// Ensure the channel is open
		if ch == nil || closeChan == nil {
			log.Println("RabbitMQ channel is closed. Reconnecting...")

			// Reconnect with retry mechanism
			for {
				err := manager.Reconnect()
				if err == nil {
					ch = manager.Channel
					closeChan = ch.NotifyClose(make(chan *amqp.Error))
					log.Println("Reconnected to RabbitMQ.")
					break
				}
				log.Printf("Failed to reconnect to RabbitMQ: %v. Retrying in 5 seconds...\n", err)
				time.Sleep(5 * time.Second)
			}
		}

		// Declare the queue
		_, err := ch.QueueDeclare(
			"image_processing",
			true,  // Durable
			false, // Auto-delete
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)
		if err != nil {
			log.Println("Failed to declare queue:", err)
			time.Sleep(5 * time.Second) // Prevent busy looping on error
			continue
		}

		// Consume messages
		msgs, err := ch.Consume(
			"image_processing",
			"",    // Consumer
			true,  // Auto-ack
			false, // Exclusive
			false, // No-local
			false, // No-wait
			nil,   // Args
		)
		if err != nil {
			log.Println("Failed to consume messages:", err)
			time.Sleep(5 * time.Second) // Prevent busy looping on error
			continue
		}

		// Process messages in a separate goroutine
		go func() {
			for msg := range msgs {
				var imageURLs []string
				if err := json.Unmarshal(msg.Body, &imageURLs); err != nil {
					log.Println("Failed to unmarshal message:", err)
					continue
				}

				// Perform image processing logic
				log.Println("Processing images:", imageURLs)

				// Example: Save data to DB (replace with actual logic)
				// db.Save(&imageURLs)
			}
		}()

		// Block until the channel is closed
		err = <-closeChan
		if err != nil {
			log.Println("RabbitMQ channel closed with error:", err)
		}
	}
}
