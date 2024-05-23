package main

import (
	"backend/connect"
	"backend/global"
	"backend/handlers"
	"backend/routes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	json "github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberi18n"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/text/language"

	firebaseAdmin "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	// "firebase.google.com/go/v4/auth"
	// "firebase.google.com/go/v4/messaging"
)

//var Conf *configs.Config

var (
	appCheck *appcheck.Client
)

func init() {
	global.LoadConfig(".")

	app, err := firebaseAdmin.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	appCheck, err = app.AppCheck(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

}

func main() {

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	micro := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(recover.New())
	micro.Use(recover.New())

	app.Mount(fmt.Sprintf("/api/%s", global.Conf.ApiVersion), micro)

	app.Use(
		logger.New(logger.Config{
			// For more options, see the Config section
			Format: "${status} - ${method} ${path}\n",
		}),
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
		logger.New(logger.Config{
			Format: "${status} - ${method} ${path}\n",
		}),
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
