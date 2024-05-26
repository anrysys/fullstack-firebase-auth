package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	// Create a new fiber app
	// app := fiber.New()

	// Run the main function in a goroutine
	go main()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Send a GET request to the server
	resp, err := http.Get("http://localhost:3000")
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
