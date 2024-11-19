package main

import (
	"net/http"
	

	"github.com/gorilla/mux"
	"product-management/config"
	"product-management/internal/api"
	"product-management/internal/cache"
	"product-management/internal/db"
	"product-management/internal/logging"
	"product-management/internal/processing"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db.InitDB()

	// Initialize Redis cache
	cache.InitRedis()

	// Initialize logger
	logger := logging.InitLogger()

	// Start the image processing consumer as a goroutine
	go func() {
		logger.Info("Starting image processing consumer...")
		processing.ConsumeQueue("product_images")
	}()

	// Setup API routes
	router := mux.NewRouter()
	router.HandleFunc("/products", api.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id:[0-9]+}", api.GetProduct).Methods("GET")
	router.HandleFunc("/products", api.GetProducts).Methods("GET")

	// Start HTTP server
	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
