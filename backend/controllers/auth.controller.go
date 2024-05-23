package controllers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"backend/connect"
	"backend/global"
	"backend/middleware"

	"backend/handlers"
	"backend/models"
	"backend/utils/messenger"

	"github.com/gofiber/contrib/fiberi18n"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "gorm.io/hints"
)

// Step 1 for pre-registration or authentication, sending a verification code by email
// Description: checking otp_code for exceeding the limit and time. Sending a message to the user. Saving the next attempt (otp_code, email) in the database.
func PreRegister(c *fiber.Ctx) error {

	// Подключаемся к базе данных
	db := connect.GetDatabase()
	var payload *models.PreRegister

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}
	var notExist bool
	var user models.User
	email := strings.ToLower(payload.Email)

	err := db.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {

		// If a new user
		if err == gorm.ErrRecordNotFound {
			notExist = true
		} else {
			// Exist user
			notExist = false
		}
	}
	// Forbidden Danger stauses for user
	// TODO Добавить {%s} в .env будет удален через {%s} месяц
	if user.UserStatus != models.UserStatusActive {
		switch user.UserStatus {
		case models.UserStatusBlocked:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsBlocked", c)})
		case models.UserStatusDeleted:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsDeleted", c)})
		//case models.UserStatusPending:
		//return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsPending", c)})
		case models.UserStatusRejected:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsRejected", c)})
		}
	}

	// Forbidden role <admin> for user
	if user.UserRole == models.UserRoleAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("HttpStatusForbidden", c)})
	}
	// Check code verification (3 attempt in the last 15 minutes)
	count, err := CountAttempts(email)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	if count >= 3 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("ErrorAttemptExceeded", c)})
	}

	// var user_status, user_role,
	var customer string
	user_status := string(user.UserStatus)
	user_role := string(user.UserRole)
	// If a new user
	if notExist {
		customer = "New Customer"
	} else {
		// Exist user
		if user.FirstName == "" {
			customer = "Customer"
		} else {
			customer = user.FirstName
		}
	}

	countInt := 5
	//otp_code := randomString(countInt)
	otp_code := fmt.Sprintf("%0"+strconv.Itoa(countInt)+"d", rand.Intn(int(math.Pow10(countInt))))

	LangSubject := fiberi18n.MustGetMessage(&i18n.LocalizeConfig{
		MessageID: "SubjectOtpCode",
		TemplateData: map[string]string{
			"OtpCode": otp_code,
		},
	})

	LangMessage := fiberi18n.MustGetMessage(&i18n.LocalizeConfig{
		MessageID: "MessageOtpCode",
		TemplateData: map[string]string{
			"Customer": customer,
			"OtpCode":  otp_code,
		},
	})

	subject := fmt.Sprint(LangSubject)
	message := fmt.Sprint(LangMessage)
	emails := strings.Split(payload.Email, ",")

	// send him the code by email
	err = messenger.SendEmail(emails, subject, message)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	} else {

		// +TODO Create record with attempt
		newAttempt := models.AuthLoginAttempt{
			Email:   &email,
			OtpCode: &otp_code,
		}

		result := db.Create(&newAttempt)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": result.Error})
		}

		// For new user
		//msg := "Code successfully addeded!"
		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{
				"status": "success",
				"msg":    handlers.L("CodeSuccessfullyAddeded", c),
				"data": fiber.Map{
					"user_status": user_status,
					"user_role":   user_role,
					"email":       &payload.Email}})
		// "otp_code":    &newAttempt.OtpCode
	}
}

// Step 2 check validate otp-code
func PreRegisterValidator(c *fiber.Ctx) error {
	db := connect.GetDatabase()
	var payload *models.PreRegisterValidator

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	email := strings.ToLower(payload.Email)
	otp_code := string(payload.OtpCode)
	lang := strings.ToLower(payload.Lang)

	// Check code verification (3 attempt in the last 15 minutes)
	count, err := CountAttempts(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	if count >= 4 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("ErrorAttemptExceeded", c)})
	}

	// Check if a record with the email and otp_code exists in auth_login_attempts.
	// If not, return an error message
	var attempt models.AuthLoginAttempt
	err = db.Where("email = ? AND otp_code = ?", email, otp_code).First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("InvalidOtpCode", c)})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
		}
	}
	passwordString := randomString(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	password := string(hashedPassword) // Convert []byte to string
	// Create or give data user
	var user models.User
	short_name := randomString(0)
	db.Where(models.User{Email: &email}).Assign(
		models.User{
			UserStatus: models.UserStatusPending,
			ShortName:  short_name,
			Password:   password,
			Lang:       lang}).FirstOrCreate(&user)

	// Forbidden Danger stauses for user
	// TODO Добавить {%s} в .env будет удален через {%s} месяц
	if user.UserStatus != models.UserStatusActive {
		switch user.UserStatus {
		case models.UserStatusBlocked:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsBlocked", c)})
		case models.UserStatusDeleted:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsDeleted", c)})
		//case models.UserStatusPending:
		//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsPending", c)})
		case models.UserStatusRejected:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("AccountIsRejected", c)})
		}
	}

	// Delete all previos attempts for user
	err = db.Where("email = ?", email).Delete(&models.AuthLoginAttempt{}).Error
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}
	// Create token
	return CreateTokenForUser(c, user)
}

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
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneCode:   payload.PhoneCode,
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": message})
	}
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
	if errAccess != nil {
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

// Check the auth_login_attempts table for the number of records with the email in the last 15 minutes.
// If there are more than 3 records, return an error with the message.
func CountAttempts(email string) (int64, error) {
	db := connect.GetDatabase()
	var count int64
	err := db.Model(&models.AuthLoginAttempt{}).Select("email").Where("email = ? AND created_at > CURRENT_TIMESTAMP - INTERVAL '15 minutes'", email).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
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
}

func Logout(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")
	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": handlers.L("TokenIsInvalidOrSessionHasExpired", c)})
	}
	ctx := context.TODO()

	// Print refresh_token claims
	fmt.Println(refresh_token)

	tokenClaims, err := middleware.ValidateToken(refresh_token, global.Conf.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	access_token_uuid := c.Locals("access_token_uuid").(string)
	_, err = connect.RedisClient.Del(ctx, tokenClaims.TokenUuid, access_token_uuid).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": handlers.L("LogoutSuccessful", c)})
}

// Random generates a random string.
func randomString(len int) string {
	if len == 0 {
		return uuid.NewString()
	}
	return uuid.NewString()[:len]
}
