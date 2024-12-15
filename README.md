## Zocket-ProductManagement
# Project Structure Overview
Here’s a breakdown of the key components in the project:

# Main Application (main.go):

Initializes and connects to the database, Redis, and RabbitMQ.
Registers product routes using the Gin framework.
Starts background consumers (like the image processing task) for handling asynchronous tasks.
Handles graceful shutdown upon termination signals.
# Configuration (config/):

InitDB(): Initializes a connection to a PostgreSQL database using GORM.
InitRedis(): Sets up a Redis client for caching.
InitRabbitMQ(): Establishes a connection to RabbitMQ, and manages reconnection logic.
# RabbitMQ Manager (config/rabbitmq_manager.go):

RabbitMQManager: Manages connections and channels to RabbitMQ.
Reconnect(): Reconnects to RabbitMQ if the connection or channel is lost, with retry logic.
# Handlers (handlers/):

RegisterProductRoutes(): Registers routes for product creation and retrieval.
Product Handlers: Handle HTTP requests for creating, retrieving, and listing products.
# Services (services/):

CreateProduct(): Saves a product to the database and publishes a message to RabbitMQ.
GetProductByID(): Retrieves a product by ID from Redis (cache) or falls back to the database.
GetProducts(): Fetches products with optional filters from the database.
PublishToQueue(): Publishes messages to RabbitMQ.
# Models (models/):

Defines the structure of User and Product models for use with the database.
Includes fields for product information like name, description, images, and price.
# Queue (queue/):

StartImageProcessingConsumer(): Consumes image processing tasks from a RabbitMQ queue and processes them asynchronously.
## Architectural Choices
1. Gin Framework:
Why Gin?: Gin is chosen for its speed and efficiency. It is a lightweight web framework for Go that provides a simple interface for routing, middleware support, and error handling, making it a suitable choice for building REST APIs.
Scalability: Gin’s performance under load allows for handling high traffic and concurrent requests effectively.
2. PostgreSQL (via GORM):
Why PostgreSQL?: PostgreSQL is a powerful relational database that supports complex queries, transactions, and indexing. GORM is used as the ORM for interacting with the database, providing an abstract layer to perform database operations with ease.
Auto-migration: GORM’s auto-migration feature ensures that the database schema stays up-to-date with the Go models.
3. Redis:
Why Redis?: Redis is used as an in-memory caching layer to store frequently accessed data like product information. By caching data in Redis, we reduce the load on the database and improve the performance of read-heavy endpoints.
Cache fallback: Redis is checked first before querying the database for product details, ensuring faster responses.
4. RabbitMQ (via amqp):
Why RabbitMQ?: RabbitMQ is used for asynchronous messaging, allowing the system to offload long-running tasks (like image processing) to background consumers. This decouples the task from the HTTP request-response cycle and improves system responsiveness.
Resiliency and Reconnection Logic: The RabbitMQ connection and channel are managed with automatic reconnection logic in case of failures, ensuring that the system remains functional under various failure scenarios.
## Setup Instructions
Prerequisites:

Install Go: Ensure Go 1.18 or higher is installed on your system.
Install PostgreSQL: Ensure PostgreSQL is installed and running.
Install Redis: Ensure Redis is installed and running.
Install RabbitMQ: Ensure RabbitMQ is installed and running.
Clone the Repository:
git clone <repository-url>
cd <project-directory>
Install Dependencies: Make sure you have the necessary Go dependencies. Run the following command:

go mod tidy
Configure Database and Redis:

Modify the dsn string in config/config.go to match your PostgreSQL credentials:
dsn := "host=localhost user=postgres password=password dbname=productdb port=5432 sslmode=disable"
Ensure that Redis and RabbitMQ are running on their respective default ports (localhost:6379 for Redis and localhost:5672 for RabbitMQ).
Run the Application:
go run main.go
This will start the Gin web server on localhost:8080.

Accessing Endpoints:

POST /products: Create a new product.
GET /products/:id: Retrieve a product by ID.
GET /products: List products with optional filters (user_id, min_price, max_price).
## Assumptions
Database: The PostgreSQL database is assumed to be running locally on the default port (5432). It is also assumed that a database named productdb exists.
Redis: The Redis instance is running locally on the default port (6379) without a password.
RabbitMQ: RabbitMQ is running locally on the default port (5672), and the default credentials (guest:guest) are used for authentication.
Message Processing: The background image processing task is simulated by logging the image URLs received from RabbitMQ; this can be extended to perform actual image processing tasks.
Environment Variables: The application does not use environment variables for configuration (like Redis or RabbitMQ credentials), but this can be added for production environments.
