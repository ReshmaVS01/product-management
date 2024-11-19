package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "product-management/internal/api"
    "product-management/internal/db"
)

func TestCreateProduct(t *testing.T) {
    payload := map[string]interface{}{
        "user_id":              1,
        "product_name":         "Test Product",
        "product_description":  "This is a test product.",
        "product_images":       []string{"http://example.com/image1.jpg"},
        "product_price":        100.00,
    }

    payloadBytes, _ := json.Marshal(payload)
    req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payloadBytes))
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(api.CreateProduct)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusCreated {
        t.Errorf("Expected status code %v, got %v", http.StatusCreated, rr.Code)
    }

    var response api.Product
    if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
        t.Fatal("Failed to parse response:", err)
    }

    if response.ProductName != "Test Product" {
        t.Errorf("Expected product name to be 'Test Product', got '%v'", response.ProductName)
    }
}
