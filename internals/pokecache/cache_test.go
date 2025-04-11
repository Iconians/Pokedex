package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(2 * time.Second)

	key := "test-key"
	val := []byte("test-data")
	cache.Add(key, val)

	cacheVal, found := cache.Get(key)
	if !found {
		t.Fatalf("Expected entry to be found in cache")
	}
	if string(cacheVal) != "test-data" {
		t.Fatalf("Expected '%s', got '%s'", val, cacheVal)
	}

	time.Sleep(3 * time.Second)

	_, found = cache.Get(key)
	if found {
		t.Fatalf("Expected entry to be reaped from cache")
	}
}
