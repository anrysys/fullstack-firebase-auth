package main

import (
	"backend/connect"
	"backend/global"
	"backend/handlers"
	"backend/routes"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	json "github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberi18n"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/text/language"
)

var (
	appCheckClient *appcheck.Client
)

func init() {

	global.LoadConfig()

}

func main() {

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${status} - ${method} ${path}\n",
	}))

	// Firebase Admin configuration
	conf := &firebase.Config{
		DatabaseURL:   global.Conf.FirebaseDatabaseURL,
		ProjectID:     global.Conf.FirebaseProjectID,
		StorageBucket: global.Conf.FirebaseStorageBucket,
	}
	// Initialize Firebase App with App Check
	firebaseApp, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialize App Check client
	appCheckClient, err = firebaseApp.AppCheck(context.Background())
	if err != nil {
		log.Fatalf("error initializing app check: %v\n", err)
	}

	// Middleware to check App Check token
	app.Use(func(c *fiber.Ctx) error {

		appCheckToken := c.Get("X-Firebase-AppCheck")

		// println(appCheckToken)

		if appCheckToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "errors": "Unauthorized - No App Check token provided"})
		}

		_, err := appCheckClient.VerifyToken(appCheckToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "errors": "Unauthorized - Invalid App Check token"})
		}

		return c.Next()
	})

	micro := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(recover.New())
	micro.Use(recover.New())

	app.Mount(fmt.Sprintf("/api/%s", global.Conf.ApiVersion), micro)

	app.Use(
		// 3 requests per 10 seconds max
		limiter.New(limiter.Config{
			Expiration: 10 * time.Second,
			Max:        3,
		}),
		cors.New(cors.Config{
			AllowOrigins:     global.Conf.ClientOrigin,
			AllowHeaders:     "Origin, Content-Type, Accept",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowCredentials: true,
		}),

		fiberi18n.New(&fiberi18n.Config{
			RootPath:         "./resources/langs",
			AcceptLanguages:  []language.Tag{language.English, language.Ukrainian, language.Russian},
			DefaultLanguage:  language.English,
			UnmarshalFunc:    toml.Unmarshal,
			FormatBundleFile: "toml",
			LangHandler:      handlers.SetLocale,
		}))

	micro.Use(
		// 3 requests per 10 seconds max
		limiter.New(limiter.Config{
			Expiration: 10 * time.Second,
			Max:        3,
		}),
		cors.New(cors.Config{
			AllowOrigins:     global.Conf.ClientOrigin,
			AllowHeaders:     "Origin, Content-Type, Accept",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowCredentials: true,
		}),

		fiberi18n.New(&fiberi18n.Config{
			RootPath:         "./resources/langs",
			AcceptLanguages:  []language.Tag{language.English, language.Ukrainian, language.Russian},
			DefaultLanguage:  language.English,
			UnmarshalFunc:    toml.Unmarshal,
			FormatBundleFile: "toml",
			LangHandler:      handlers.SetLocale,
		}))

	// Подключаемся к базе данных
	connect.GetDatabase()

	// Подключаемся к Redis
	connect.ConnectRedis()

	// Инициализируем маршруты
	routes.Setup(micro)

	// Запускаем сервер
	port := fmt.Sprintf(":%s", global.Conf.ServerPort)
	log.Fatal(app.Listen(port))

}
