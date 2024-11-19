**🚀 Product Management System with Asynchronous Image Processing**
Welcome to the Product Management System, a backend application built using Golang that focuses on asynchronous image processing, caching, and high scalability. This system enables efficient management of products while processing and compressing images in the background.

**🛠 Features**
API Design:
POST /products: Add a new product with details including images.
GET /products/{id}: Retrieve product details by ID, with processed image data.
GET /products: Fetch products with filters (user ID, price range, name).
Asynchronous Image Processing:
Compress and store images using RabbitMQ for message queuing.
Upload compressed images to S3-like storage.
Data Storage:
PostgreSQL database for persistent storage.
Store original and compressed image URLs in separate columns.
Caching:
Use Redis to cache product details for faster retrieval.
Implement cache invalidation to reflect real-time updates.
Enhanced Logging:
Structured logging with logrus for debugging and monitoring.
Error Handling:
Robust error handling and retry mechanisms.
Testing:
Comprehensive unit and integration tests with 90%+ coverage.

**🏗️ System Architecture**
Overview
Modular architecture ensures scalability and maintainability.
Asynchronous image processing decouples tasks for better performance.
Redis caching reduces database load and improves response times.
Key Components
RESTful APIs:
Built using net/http and gorilla/mux.
Database:
PostgreSQL with GORM ORM.
Caching:
Redis for storing frequently accessed data.
Message Queue:
RabbitMQ for task distribution.
Image Processing:
disintegration/imaging library for compression.

**🗂️ Directory Structure**
product-management/
├── cmd/
│   └── main.go                # Entry point of the application
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── api/                   # API handlers
│   │   └── products.go
│   ├── cache/                 # Redis caching
│   │   └── cache.go
│   ├── db/                    # Database connection and setup
│   │   └── db.go
│   ├── logging/               # Structured logging
│   │   └── logger.go
│   ├── processing/            # Image processing consumer
│   │   └── consumer.go
│   └── queue/                 # Message queue publisher
│       └── publisher.go
├── tests/
│   ├── api_test.go            # Unit tests for API endpoints
│   ├── cache_test.go          # Unit tests for caching functionality
│   ├── integration_test.go    # Integration tests for end-to-end validation
│   └── benchmark_test.go      # Benchmark tests for performance analysis
├── .gitignore                 # Files and directories to exclude from version control
└── README.md                  # Project documentation


**🚀 Getting Started**
Prerequisites
Golang (v1.18+)
PostgreSQL
Redis
RabbitMQ

**Setup Instructions**
Clone the Repository:
code:
git clone https://github.com/ReshmaVS01/product-management.git
cd product-management

Set Up Environment Variables: Create a .env file in the root directory with the following:

code:
POSTGRES_HOST=localhost
POSTGRES_USER=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_DB=product_management
POSTGRES_PORT=5432

Run PostgreSQL and Redis:

Ensure PostgreSQL and Redis are running (I used brew services) locally or use Docker:
code:
brew services start redis
brew services start rabbitmq
or
docker-compose up

Install Dependencies:
code:
go mod tidy

Run the Application:
code:
go run cmd/main.go

**🧪 Testing**
Run the tests to ensure everything is functioning correctly:

code:
go test ./tests/...

**📊 Performance**
Optimized for high performance using asynchronous processing.
Redis caching ensures low-latency responses for frequently accessed data.
Logs every request, processing status, and failures.

**📚 Documentation**
API Endpoints:
Full details available in internal/api/.
Database Schema:
Automatically created using GORM models.

**💻 Tech Stack**
Programming Language: Golang
Database: PostgreSQL
Cache: Redis
Queue: RabbitMQ
Logging: Logrus
Image Processing: Imaging library

