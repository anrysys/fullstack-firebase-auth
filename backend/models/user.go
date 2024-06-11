package models

import (
	"database/sql"
	"encoding/json"
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
	UserRoleAdmin    UserRoles = "admin"
	UserRoleCustomer UserRoles = "customer"
	UserRoleSupplier UserRoles = "supplier"
	UserRoleCourier  UserRoles = "courier"
)

type User struct {
	// gorm.Model
	ID           *int64       `gorm:"column:user_id;not null;primaryKey" json:"user_id"`
	Email        *string      `gorm:"type:varchar(100);not null;uniqueIndex:users_email_uidx"`
	UserSocialID int64        `gorm:"column:user_social_id;"`
	FirebaseUid  string       `gorm:"type:varchar(100);not null;"`
	ProviderId   string       `gorm:"type:varchar(100);not null;"`
	ShortName    string       `gorm:"type:varchar(36)"`
	FirstName    string       `gorm:"type:varchar(50)"`
	LastName     string       `gorm:"type:varchar(50)"`
	PhoneNumber  string       `gorm:"type:varchar(20)"`
	Password     string       `gorm:"type:varchar(100)"`
	Lang         string       `gorm:"type:varchar(2)"`
	UserStatus   UserStatuses `gorm:"not null;default:pending"`
	UserRole     UserRoles    `gorm:"not null;default:customer"`
	Photo        string       `gorm:"type:varchar(100)'"`
	ActivatedAt  sql.NullTime `gorm:""`
	VerifiedAt   sql.NullTime `gorm:""`
	// CreatedAt time.Time `gorm:"default:current_timestamp"`
	CreatedAt *time.Time   `gorm:""`
	UpdatedAt sql.NullTime `gorm:""`
	DeletedAt sql.NullTime `gorm:"index"` // Soft delete
}

type UserSocial struct {
	// gorm.Model
	ID            *int64          `gorm:"column:user_social_id;not null;primaryKey" json:"user_social_id"`
	UserID        *int64          `gorm:"column:user_id;not null;" json:"user_id"`
	FirebaseUid   string          `gorm:"column:firebase_uid;type:varchar(100);not null;uniqueIndex:users_firebase_uid_uidx"`
	ProviderId    string          `gorm:"type:varchar(100);"`
	Email         *string         `gorm:"type:varchar(100);"`
	TenantId      *string         `gorm:"type:varchar(100);"`
	EmailVerified bool            `gorm:"not null;default:false"`
	PhoneNumber   string          `gorm:"type:varchar(20)"`
	IsAnonymous   bool            `gorm:"not null;default:false"`
	Username      string          `gorm:"type:varchar(36)"`
	PhotoUrl      string          `gorm:"type:varchar(150)"`
	ProviderData  json.RawMessage `gorm:"type:json"`
	CreatedAt     *time.Time      `gorm:""`
	UpdatedAt     sql.NullTime    `gorm:""`
	DeletedAt     sql.NullTime    `gorm:""`
}

// uid: 'LnYv4WD6jbVxHlArKrj4AhzPqaF3',
// providerId: 'github.com',
// tenantId: null,
// isNewUser: false,
// isAnonymous: false,
// displayName: 'anrysys',
// email: 'anrysys@gmail.com',
// emailVerified: false,
// phoneNumber: null,
// photoURL: 'https://avatars.githubusercontent.com/u/1523609?v=4',
// lang: 'ru'

type UserSocialRequest struct {
	FirebaseUid   string          `json:"uid" validate:"required"`
	ProviderId    string          `json:"providerId" validate:"required"`
	Lang          string          `json:"lang" validate:"required,oneof=en uk es fr de it pt zh ja ko ru"`
	Email         string          `json:"email" validate:"omitempty,email"`
	PhoneNumber   string          `json:"phoneNumber"`
	Username      string          `json:"displayName"`
	PhotoUrl      string          `json:"photoURL"`
	EmailVerified bool            `json:"emailVerified"`
	IsAnonymous   bool            `json:"isAnonymous"`
	TenantId      string          `json:"tenantId"`
	IsNewUser     bool            `json:"isNewUser"`
	ProviderData  json.RawMessage `json:"providerData"`
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

type ResetPassword struct {
	Lang  string `json:"lang" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID          int64        `json:"user_id,omitempty"`
	ShortName   string       `json:"short_name,omitempty"`
	FirstName   string       `json:"first_name,omitempty"`
	LastName    string       `json:"last_name,omitempty"`
	PhoneNumber string       `json:"phone_number,omitempty"`
	Email       string       `json:"email,omitempty"`
	Photo       string       `json:"photo,omitempty"`
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
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Photo:       user.Photo,
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
