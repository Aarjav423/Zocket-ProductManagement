package handlers

import (
	"log"
	"net/http"
	"project/models"
	"project/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

// RegisterProductRoutes registers product routes.
func RegisterProductRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client, rabbitCh *amqp.Channel) {
	r.POST("/products", createProductHandler(db))
	r.GET("/products/:id", getProductByIDHandler(db, redisClient))
	r.GET("/products", getProductsHandler(db))
}

func createProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product

		// Bind JSON request to the product model
		if err := c.ShouldBindJSON(&product); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Debug: Check the received product
		log.Printf("Received product: %+v", product)

		// Insert product into the database
		if err := db.Create(&product).Error; err != nil {
			log.Println("Error inserting product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		// Debug: Check the product after saving (to see if fields are filled)
		log.Printf("Product saved: %+v", product)

		// Respond with the created product
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": product})
	}
}

// getProductByIDHandler handles fetching a product by ID
func getProductByIDHandler(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Fetch product from Redis or fallback to DB
		product, err := services.GetProductByID(db, redisClient, id)
		if err != nil {
			log.Println("Error fetching product:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

// getProductsHandler handles fetching products with filters
func getProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query products with filters
		products, err := services.GetProducts(db, c.Query("user_id"), c.Query("min_price"), c.Query("max_price"))
		if err != nil {
			log.Println("Error fetching products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
