package routes

import (
	"backend/connect"
	"backend/controllers"
	"backend/middleware"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func Setup(micro *fiber.App) {

	ctx := context.TODO()
	value, err := connect.RedisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	micro.Route("/auth", func(router fiber.Router) {
		// router.Post("/preregister", controllers.PreRegister)
		// router.Post("/preregistervalidator", controllers.PreRegisterValidator)
		router.Post("/register", controllers.Register)
		router.Post("/login", controllers.Login)
		router.Get("/logout", middleware.Auth, controllers.Logout)
		router.Get("/refresh", controllers.RefreshAccessToken)
		router.Post("/login-socials", controllers.LoginSocials)
	})

	micro.Get("/users/me", middleware.Auth, controllers.GetMe)

	// micro.Get("/users/me", middleware.FirebaseAppCheck(func(c *fiber.Ctx) error {
	// 	return c.Next()
	// }), controllers.GetMe)

	// OR you can use this way
	//micro.Get("/users/me", middleware.FirebaseAppCheck(middleware.Auth), controllers.GetMe)

	micro.Patch("/users/updateme", middleware.Auth, controllers.UpdateMe)
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": value,
		})
	})

	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

}
