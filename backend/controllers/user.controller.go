package controllers

import (
	"database/sql"
	"time"

	"backend/connect"

	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

// Step 3 (post register) update data for user (first name and last name)
func UpdateMe(c *fiber.Ctx) error {
	db := connect.GetDatabase()

	var payload *models.UpdateMe

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	first_name := string(payload.FirstName)
	last_name := string(payload.LastName)
	user := c.Locals("user").(models.UserResponse)

	// TODO Добавить {%s} в .env будет удален через {%s} месяц
	if user.UserStatus != string(models.UserStatusActive) {
		switch user.UserStatus {
		case string(models.UserStatusBlocked):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": utils.L("AccountIsBlocked", c)})
		case string(models.UserStatusDeleted):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": utils.L("AccountIsDeleted", c)})
		case string(models.UserStatusRejected):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": utils.L("AccountIsRejected", c)})
		}
	}
	var activatedAt = sql.NullTime{}
	// USE ONLY ONE TIME! If already activated - STOP! Only for use on step 3 !!! Resctrict for ordinary condition !!!
	if user.ActivatedAt != activatedAt {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"status": "fail", "errors": utils.L("ServerStatusMethodNotAllowed", c)})
	}
	var msg = utils.L("SavedSuccessfully", c)
	if user.UserStatus == string(models.UserStatusPending) {
		msg = utils.L("RegisterUserConglaturation", c)
		activatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	u := models.User{
		FirstName:   first_name,
		LastName:    last_name,
		UserStatus:  models.UserStatusActive,
		ActivatedAt: activatedAt,
	}

	db.Where(models.User{ID: &user.ID}).Updates(&u)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": msg})

}
