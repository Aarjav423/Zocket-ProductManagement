package main

import (
	"log"
	"os"
	"os/signal"
	"project/config"
	"project/handlers"
	"project/queue"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configurations
	db := config.InitDB()                         // Initialize the database (GORM)
	redisClient := config.InitRedis()             // Initialize Redis client
	rabbitManager := config.InitRabbitMQManager() // Initialize RabbitMQ connection and channel

	// Ensure RabbitMQ resources are cleaned up on exit
	defer func() {
		if rabbitManager.Connection != nil {
			if err := rabbitManager.Connection.Close(); err != nil {
				log.Printf("Failed to close RabbitMQ connection: %v", err)
			}
		}
		if rabbitManager.Channel != nil {
			if err := rabbitManager.Channel.Close(); err != nil {
				log.Printf("Failed to close RabbitMQ channel: %v", err)
			}
		}
	}()

	// Start background consumers (such as image processing)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Consumer panic recovered: %v", r)
			}
		}()
		// Start the image processing consumer
		queue.StartImageProcessingConsumer(db, rabbitManager)
	}()

	// Initialize the Gin router
	router := gin.Default()

	// Register product routes with the Gin router
	handlers.RegisterProductRoutes(router, db, redisClient, rabbitManager.Channel)

	// Graceful shutdown handling
	gracefulShutdown(router, rabbitManager)
}

// gracefulShutdown handles cleanup when the server is terminated
func gracefulShutdown(router *gin.Engine, rabbitManager *config.RabbitMQManager) {
	// Create a channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		log.Println("Server running on :8080")
		if err := router.Run(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for a termination signal (SIGINT or SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Cleanup RabbitMQ resources
	if rabbitManager.Connection != nil {
		if err := rabbitManager.Connection.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}

	if rabbitManager.Channel != nil {
		if err := rabbitManager.Channel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}

	// Shutdown the router gracefully (optional cleanup)
	log.Println("Server shutdown complete.")
}
