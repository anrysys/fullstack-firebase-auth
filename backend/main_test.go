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

// func TestVerifyToken(t *testing.T) {
//     ctx := context.Background()
//     opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
//     app, err := firebase.NewApp(ctx, nil, opt)
//     if err != nil {
//         t.Fatalf("error initializing app: %v", err)
//     }

//     client, err := app.AppCheck(ctx)
//     if err != nil {
//         t.Fatalf("error getting AppCheck client: %v", err)
//     }

//     // Replace this with a valid App Check token.
//     token := "REPLACE_THIS_WITH_A_VALID_TOKEN"

//     _, err = client.VerifyToken(ctx, token)
//     if err != nil {
//         t.Fatalf("VerifyToken failed: %v", err)
//     }
// }
