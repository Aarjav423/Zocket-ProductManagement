package services

import (
	"encoding/json"
	"log"

	"project/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

// GenerateID generates a new unique ID for the product.
func GenerateID() string {
	return uuid.New().String()
}

// CreateProduct creates a product in the database and publishes a message to RabbitMQ.
func CreateProduct(db *gorm.DB, rabbitCh *amqp.Channel, product *models.Product) error {
	// Save product in the database
	if err := db.Create(product).Error; err != nil {
		return err
	}

	// Publish product data to RabbitMQ
	message, _ := json.Marshal(product)
	return PublishToQueue(rabbitCh, "product_queue", message)
}

// GetProductByID retrieves a product by its ID from Redis or DB.
func GetProductByID(db *gorm.DB, redisClient *redis.Client, id string) (*models.Product, error) {
	// Check Redis first
	var product models.Product
	val, err := redisClient.Get(redisClient.Context(), id).Result()
	if err == nil && val != "" {
		// Product found in cache
		json.Unmarshal([]byte(val), &product)
		return &product, nil
	}

	// If not in cache, check DB
	if err := db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Store product in Redis cache
	redisClient.Set(redisClient.Context(), id, val, 0)
	return &product, nil
}

// GetProducts retrieves products from the DB with optional filtering.
func GetProducts(db *gorm.DB, userID string, minPrice string, maxPrice string) ([]models.Product, error) {
	var products []models.Product
	query := db.Where("user_id = ?", userID)

	if minPrice != "" {
		query = query.Where("product_price >= ?", minPrice)
	}

	if maxPrice != "" {
		query = query.Where("product_price <= ?", maxPrice)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// PublishToQueue publishes a message to a RabbitMQ queue.
func PublishToQueue(ch *amqp.Channel, queueName string, message []byte) error {
	// Declare the queue (idempotent)
	_, err := ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Println("Error declaring queue:", err)
		return err
	}

	// Publish the message
	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Println("Error publishing message:", err)
		return err
	}

	log.Println("Message published to queue:", queueName)
	return nil
}
