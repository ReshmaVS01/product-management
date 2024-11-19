package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "product-management/internal/api"
    "product-management/internal/db"
)

func TestGetProductIntegration(t *testing.T) {
    // Insert a test product into the database
    product := api.Product{
        UserID:               1,
        ProductName:          "Integration Test Product",
        ProductDescription:   "Description",
        ProductImages:        []string{"http://example.com/test.jpg"},
        CompressedProductImages: []string{"http://example.com/compressed.jpg"},
        ProductPrice:         50.00,
    }
    if err := db.DB.Create(&product).Error; err != nil {
        t.Fatal("Failed to create test product:", err)
    }

    req, err := http.NewRequest("GET", "/products/"+string(product.ID), nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(api.GetProduct)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
    }

    var response api.Product
    if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
        t.Fatal("Failed to parse response:", err)
    }

    if response.ProductName != product.ProductName {
        t.Errorf("Expected product name '%s', got '%s'", product.ProductName, response.ProductName)
    }
}

