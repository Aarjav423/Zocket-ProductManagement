package config

import (
	"log"
	"project/models" // Import the models package

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the PostgreSQL database connection.
func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=aarjav dbname=productdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables for User and Product models from the models package
	db.AutoMigrate(&models.User{}, &models.Product{})
	return db
}

// InitRedis initializes the Redis client for caching purposes.
func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// InitRabbitMQ initializes the RabbitMQ connection and channel.
func InitRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}
	return conn, ch
}
