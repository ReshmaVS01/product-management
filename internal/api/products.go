package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"product-management/internal/cache"
	"product-management/internal/db"
	"product-management/internal/queue"
)

// Product struct represents the product entity
type Product struct {
	ID                     uint     `json:"id" gorm:"primaryKey"`
	UserID                 uint     `json:"user_id"`
	ProductName            string   `json:"product_name"`
	ProductDescription     string   `json:"product_description"`
	ProductImages          []string `json:"product_images" gorm:"type:text[]"`
	CompressedProductImages []string `json:"compressed_product_images" gorm:"type:text[]"`
	ProductPrice           float64  `json:"product_price"`
}

// POST /products - Create a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Save product in database
	if err := db.DB.Create(&product).Error; err != nil {
		http.Error(w, "Failed to save product", http.StatusInternalServerError)
		return
	}

	// Publish product_images URLs to the queue
	if err := queue.PublishMessage("product_images", product.ProductImages); err != nil {
		http.Error(w, "Failed to enqueue image processing task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GET /products/:id - Retrieve product details by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Check Redis cache
	cacheKey := "product:" + strconv.Itoa(id)
	cachedProduct, err := cache.Get(cacheKey)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedProduct))
		return
	}

	// Fetch product from database if not cached
	var product Product
	if err := db.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		}
		return
	}

	// Cache product data
	productJSON, _ := json.Marshal(product)
	cache.Set(cacheKey, string(productJSON))

	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}

// GET /products - Retrieve products with optional filters
func GetProducts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	productName := r.URL.Query().Get("product_name")

	var products []Product
	query := db.DB.Model(&Product{})

	// Apply filters
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if minPrice != "" {
		query = query.Where("product_price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("product_price <= ?", maxPrice)
	}
	if productName != "" {
		query = query.Where("product_name ILIKE ?", "%"+productName+"%")
	}

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

