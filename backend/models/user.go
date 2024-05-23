// package models

// import (
// 	"gorm.io/gorm"
// )

// // User модель пользователя для GORM
// type User struct {
// 	gorm.Model
// 	Username string `gorm:"uniqueIndex"`
// 	Password string
// 	Email    string `gorm:"uniqueIndex"`
// }

package models

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserStatuses string

const (
	UserStatusPending   UserStatuses = "pending"
	UserStatusActive    UserStatuses = "active"
	UserStatusSuspended UserStatuses = "suspended"
	UserStatusBlocked   UserStatuses = "blocked"
	UserStatusDeleted   UserStatuses = "deleted"
	UserStatusRejected  UserStatuses = "rejected"
)

type UserRoles string

const (
	UserRoleAdmin  UserRoles = "admin"
	UserRoleBuyer  UserRoles = "buyer"
	UserRoleSeller UserRoles = "seller"
)

// type:int64;
type User struct {
	//gorm.Model
	ID          *int64       `gorm:"column:user_id;not null;primaryKey" json:"user_id"`
	Email       *string      `gorm:"type:varchar(100);not null;uniqueIndex:users_email_uidx"`
	ShortName   string       `gorm:"type:varchar(36)"`
	FirstName   string       `gorm:"type:varchar(50)"`
	LastName    string       `gorm:"type:varchar(50)"`
	PhoneCode   string       `gorm:"type:varchar(10)"`
	PhoneNumber string       `gorm:"type:varchar(20)"`
	Password    string       `gorm:"type:varchar(100)"`
	OtpCode     string       `gorm:"type:char(5)"`
	Lang        string       `gorm:"type:varchar(2)"`
	UserStatus  UserStatuses `gorm:"not null;default:pending"`
	UserRole    UserRoles    `gorm:"not null;default:buyer"`
	Photo       string       `gorm:"type:varchar(100)'"`
	ActivatedAt sql.NullTime `gorm:""`
	VerifiedAt  sql.NullTime `gorm:""`
	// CreatedAt time.Time `gorm:"default:current_timestamp"`
	CreatedAt *time.Time   `gorm:""`
	UpdatedAt sql.NullTime `gorm:""`
	DeletedAt sql.NullTime `gorm:""`
}

// Registration for APP (step 1). Sending the first request by the user, sent by email.
type PreRegister struct {
	Email string `json:"email" validate:"required,email"`
	Lang  string `json:"lang" validate:"required"`
}

// Registration for APP (step 2). User confirmation of the code received by email.
type PreRegisterValidator struct {
	Email   string `json:"email" validate:"required,email"`
	OtpCode string `json:"otp_code" validate:"required"`
	Lang    string `json:"lang" validate:"required"`
}

// Data Additional by Registration for APP (step 3). User data (first name, last name)
type UpdateMe struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Lang      string `json:"lang" validate:"required"`
}

// Registration for administartors (ERM/CRM)
type Register struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
	Lang            string `json:"lang" validate:"required"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	PhoneCode       string `json:"phone_code,omitempty"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	Photo           string `json:"photo,omitempty"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Lang     string `json:"lang" validate:"required"`
}

// type Logout struct {
// 	Lang string `json:"lang" validate:"required"`
// }

// type Refresh struct {
// 	Lang string `json:"lang" validate:"required"`
// }

type UserResponse struct {
	ID          int64        `json:"user_id,omitempty"`
	ShortName   string       `json:"short_name,omitempty"`
	FirstName   string       `json:"first_name,omitempty"`
	LastName    string       `json:"last_name,omitempty"`
	PhoneCode   string       `json:"phone_code,omitempty"`
	PhoneNumber string       `json:"phone_number,omitempty"`
	Email       string       `json:"email,omitempty"`
	Photo       string       `json:"photo,omitempty"`
	OtpCode     string       `json:"otp_code,omitempty"`
	UserStatus  string       `json:"user_status,omitempty"`
	UserRole    string       `json:"user_role,omitempty"`
	Lang        string       `json:"lang,omitempty"`
	ActivatedAt sql.NullTime `json:"activated_at"`
	VerifiedAt  sql.NullTime `json:"verified_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:          *user.ID,
		Email:       *user.Email,
		ShortName:   user.ShortName,
		PhoneCode:   user.PhoneCode,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Photo:       user.Photo,
		OtpCode:     user.OtpCode,
		Lang:        user.Lang,
		UserStatus:  string(user.UserStatus),
		UserRole:    string(user.UserRole),
		ActivatedAt: user.ActivatedAt,
		VerifiedAt:  user.VerifiedAt,
		CreatedAt:   *user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
