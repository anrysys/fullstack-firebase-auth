package controllers

import (
	"context"
	"strings"
	"time"

	"backend/connect"
	"backend/global"
	"backend/middleware"

	"backend/handlers"
	"backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "gorm.io/hints"
)

// Register new customer
func Register(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	var payload *models.Register

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})

	}
	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": handlers.L("PasswordDoNotMatch", c)})

	}
	// Password Generic
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	password := string(hashedPassword) // Convert []byte to string

	email := strings.ToLower(payload.Email)

	short_name := randomString(0)
	newUser := models.User{
		ShortName:   short_name,
		Email:       &email,
		UserStatus:  models.UserStatusActive,
		UserRole:    models.UserRoleCustomer,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
		Password:    password,
		Photo:       payload.Photo,
		Lang:        payload.Lang,
	}

	result := db.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "errors": handlers.L("UserEmailExist", c)})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": handlers.L("SometingWentWrongServer", c)})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": models.FilterUserRecord(&newUser)}})
}

// Authentification for customer
func Login(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	var payload *models.Login

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	message := handlers.L("InvalidEmailOrPassword", c)
	var user models.User
	err := db.First(&user, "email = ?", strings.ToLower(payload.Email)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": message})
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
		}
	}
	// Disable password check
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	// if err != nil {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": message})
	// }
	// TODO Добавить {%s} в .env будет удален через {%s} месяц
	if user.UserStatus != models.UserStatusActive {
		switch user.UserStatus {
		case models.UserStatusBlocked:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsBlocked", c)})
		case models.UserStatusDeleted:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsDeleted", c)})
		case models.UserStatusPending:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsPending", c)})
		case models.UserStatusRejected:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsRejected", c)})
		}
	}

	// Create token
	return CreateTokenForUser(c, user)
}

func LoginSocials(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	var payload *models.UserSocialRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})

	}
	email := strings.ToLower(payload.Email)

	// If user NOT exists (by email) in table users, create new user in table users and in table user_socials
	// Create Full Schema: users -> user_socials
	var user models.User
	err := db.First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user
			newUser := models.User{
				FirebaseUid: payload.FirebaseUid,
				ProviderId:  payload.ProviderId,
				ShortName:   payload.Username,
				Email:       &payload.Email,
				UserStatus:  models.UserStatusActive,
				UserRole:    models.UserRoleCustomer,
				PhoneNumber: payload.PhoneNumber,
				Photo:       payload.PhotoUrl,
				Lang:        payload.Lang,
			}
			result := db.Create(&newUser)
			if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "errors": handlers.L("UserEmailExist", c)})
			} else if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": handlers.L("SometingWentWrongServer", c)})
			}

			newUserSocial := models.UserSocial{
				UserID:        newUser.ID,
				FirebaseUid:   payload.FirebaseUid,
				ProviderId:    payload.ProviderId,
				Email:         &payload.Email,
				TenantId:      &payload.TenantId,
				EmailVerified: payload.EmailVerified,
				PhoneNumber:   payload.PhoneNumber,
				IsAnonymous:   payload.IsAnonymous,
				Username:      payload.Username,
				PhotoUrl:      payload.PhotoUrl,
				ProviderData:  payload.ProviderData,
			}
			result = db.Create(&newUserSocial)
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": handlers.L("SometingWentWrongServer", c)})
			}

			// Update users table (columns: user_social_id, firebase_uid, provider_id)
			result = db.Model(&user).Where("user_id = ?", newUser.ID).Updates(models.User{
				UserSocialID: *newUserSocial.ID,
				FirebaseUid:  payload.FirebaseUid,
				ProviderId:   payload.ProviderId,
				ShortName:    payload.Username,
				PhoneNumber:  payload.PhoneNumber,
				Photo:        payload.PhotoUrl,
				Lang:         payload.Lang,
			})

			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": result.Error.Error()})
			}

			return CreateTokenForUser(c, newUser)
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
		}
	}

	// If user exists, check if user_socials exists
	// If the entry is missing, create a new entry in the user_socials table with the data
	// from the query and update the entry in the users table (columns: user_social_id, firebase_uid, provider_id)
	err = db.Joins("JOIN user_socials ON users.user_id = user_socials.user_id").First(&user, "users.email = ? AND user_socials.firebase_uid = ? AND user_socials.provider_id = ?", email, payload.FirebaseUid, payload.ProviderId).Error
	if err != nil {
		// If user_socials not exists, create new user_socials record with data from request
		if err == gorm.ErrRecordNotFound {
			newUserSocial := models.UserSocial{
				UserID:        user.ID,
				FirebaseUid:   payload.FirebaseUid,
				ProviderId:    payload.ProviderId,
				Email:         &payload.Email,
				TenantId:      &payload.TenantId,
				EmailVerified: payload.EmailVerified,
				PhoneNumber:   payload.PhoneNumber,
				IsAnonymous:   payload.IsAnonymous,
				Username:      payload.Username,
				PhotoUrl:      payload.PhotoUrl,
				ProviderData:  payload.ProviderData,
			}
			result := db.Create(&newUserSocial)
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": handlers.L("SometingWentWrongServer", c)})
			}

			// Update users table (columns: user_social_id, firebase_uid, provider_id)
			result = db.Model(&user).Updates(models.User{
				UserSocialID: *newUserSocial.ID,
				FirebaseUid:  payload.FirebaseUid,
				ProviderId:   payload.ProviderId,
				ShortName:    payload.Username,
				PhoneNumber:  payload.PhoneNumber,
				Photo:        payload.PhotoUrl,
				Lang:         payload.Lang,
			})
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": result.Error.Error()})
			}

		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
		}
	}

	// Create token
	return CreateTokenForUser(c, user)

}

// Logout user with a system and delete token from redis
func Logout(c *fiber.Ctx) error {
	ctx := context.TODO()
	access_token_uuid := c.Locals("access_token_uuid").(string)
	err := connect.RedisClient.Del(ctx, access_token_uuid).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	c.ClearCookie("logged_in")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// Create token for user
func CreateTokenForUser(c *fiber.Ctx, user models.User) error {

	accessTokenDetails, err := middleware.CreateToken(
		*user.ID,
		string(*user.Email),
		string(user.ShortName),
		string(user.UserStatus),
		string(user.UserRole),
		global.Conf.AccessTokenExpiresIn,
		global.Conf.AccessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	refreshTokenDetails, err := middleware.CreateToken(
		*user.ID,
		string(*user.Email),
		string(user.ShortName),
		string(user.UserStatus),
		string(user.UserRole),
		global.Conf.RefreshTokenExpiresIn,
		global.Conf.RefreshTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	ctx := context.TODO()
	now := time.Now()

	errAccess := connect.RedisClient.Set(ctx, accessTokenDetails.TokenUuid, user.ID, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": errAccess.Error()})
	}

	errRefresh := connect.RedisClient.Set(ctx, refreshTokenDetails.TokenUuid, user.ID, time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errRefresh != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": errRefresh.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   global.Conf.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   global.Conf.Host,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		MaxAge:   global.Conf.RefreshTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   global.Conf.Host,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   global.Conf.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
		Domain:   global.Conf.Host,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":        "success",
		"access_token":  accessTokenDetails.Token,
		"refresh_token": refreshTokenDetails.Token,
	})
}

func RefreshAccessToken(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	refresh_token := c.Cookies("refresh_token")
	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("CouldNotRefreshAccessToken", c)})
	}

	ctx := context.TODO()

	tokenClaims, err := middleware.ValidateToken(refresh_token, global.Conf.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	userid, err := connect.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("CouldNotRefreshAccessToken", c)})
	}

	var user models.User
	err = db.First(&user, "user_id = ?", userid).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("TheUserBelongingToThisTokenNoLoggerExists", c)})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})

		}
	}

	accessTokenDetails, err := middleware.CreateToken(
		*user.ID,
		string(*user.Email),
		string(user.ShortName),
		string(user.UserStatus),
		string(user.UserRole),
		global.Conf.AccessTokenExpiresIn,
		global.Conf.AccessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	now := time.Now()

	errAccess := connect.RedisClient.Set(ctx, accessTokenDetails.TokenUuid, user.ID, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": errAccess.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   global.Conf.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   global.Conf.Host,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   global.Conf.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
		Domain:   global.Conf.Host,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
} // Random generates a random string.
func randomString(len int) string {
	if len == 0 {
		return uuid.NewString()
	}
	return uuid.NewString()[:len]
}
