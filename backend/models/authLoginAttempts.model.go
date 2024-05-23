package models

import (
	"time"
	//"github.com/jackc/pgtype"
)

type AuthLoginAttempt struct {
	ID      *int64  `json:"auth_login_attempt_id" gorm:"column:auth_login_attempt_id;type:int64;not null;primaryKey"`
	Email   *string `gorm:"type:varchar(100);not null"`
	OtpCode *string `gorm:"type:char(5);not null"`
	//ClientIp  *pgtype.Inet `json:"client_ip" gorm:"type:inet;not null"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	//UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

type AuthLoginAttemptResponse struct {
	ID      int64  `json:"auth_login_attempt_id,omitempty"`
	Email   string `json:"email,omitempty"`
	OtpCode string `json:"otp_code,omitempty"`
	//ClientIp  pgtype.Inet `json:"client_ip,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
}

func FilterAuthLoginAttemptRecord(attempt *AuthLoginAttempt) AuthLoginAttemptResponse {
	return AuthLoginAttemptResponse{
		ID:      *attempt.ID,
		Email:   *attempt.Email,
		OtpCode: *attempt.OtpCode,
		//ClientIp:  *attempt.ClientIp,
		CreatedAt: *attempt.CreatedAt,
		//UpdatedAt: *attempt.UpdatedAt,
	}
}
