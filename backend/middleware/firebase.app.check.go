package middleware

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	// "firebase.google.com/go/v4/auth"
	// "firebase.google.com/go/v4/messaging"
)

var (
	appCheck *appcheck.Client
)

func FirebaseAppCheck(handler fiber.Handler) fiber.Handler {

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	appCheck, err = app.AppCheck(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return func(c *fiber.Ctx) error {
		appCheckToken := c.Get("X-Firebase-AppCheck")
		if appCheckToken == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized.")
		}

		_, err := appCheck.VerifyToken(appCheckToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized.")
		}

		// If VerifyToken() succeeds, continue with the provided handler.
		return handler(c)

	}
}
