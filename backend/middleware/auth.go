package middleware

import (
	"context"
	"strings"

	"backend/connect"
	"backend/global"

	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Auth(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	var access_token string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		access_token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		access_token = c.Cookies("access_token")
	}
	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "errors": utils.L("YouAreNotLoggedIn", c)})
	}

	tokenClaims, err := ValidateToken(access_token, global.Conf.AccessTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	ctx := context.TODO()
	userid, err := connect.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": utils.L("TokenIsInvalidOrSessionHasExpired", c)})
	}

	var user models.User
	err = db.First(&user, "user_id = ?", userid).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": utils.L("TheUserBelongingToThisTokenNoLoggerExists", c)})
	}

	c.Locals("user", models.FilterUserRecord(&user))
	c.Locals("access_token_uuid", tokenClaims.TokenUuid)

	return c.Next()
}
