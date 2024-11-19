package tests

import (
    "testing"
    "product-management/internal/cache"
)

func TestRedisCache(t *testing.T) {
    key := "test_key"
    value := "test_value"

    cache.Set(key, value)
    cachedValue, err := cache.Get(key)
    if err != nil {
        t.Fatalf("Failed to get value from cache: %v", err)
    }

    if cachedValue != value {
        t.Errorf("Expected value '%s', got '%s'", value, cachedValue)
    }
}

