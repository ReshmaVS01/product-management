package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "product-management/internal/api"
)

func BenchmarkGetProduct(b *testing.B) {
    req, err := http.NewRequest("GET", "/products/1", nil)
    if err != nil {
        b.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(api.GetProduct)

    for i := 0; i < b.N; i++ {
        handler.ServeHTTP(rr, req)
    }
}

