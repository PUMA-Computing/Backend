package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                     uuid.UUID  `pg:"type:uuid" json:"id"`
	Username               string     `json:"username"`
	Password               string     `json:"password"`
	FirstName              string     `json:"first_name"`
	MiddleName             *string    `json:"middle_name"`
	LastName               string     `json:"last_name"`
	Email                  string     `json:"email"`
	StudentID              string     `json:"student_id"`
	Major                  string     `json:"major"`
	ProfilePicture         string     `json:"profile_picture"`
	DateOfBirth            *time.Time `json:"date_of_birth"`
	RoleID                 int        `json:"role_id"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	Year                   string     `json:"year"`
	EmailVerified          bool       `json:"email_verified"`
	EmailVerificationToken string     `json:"email_verification_token"`
	PasswordResetToken     string     `json:"password_reset_token"`
	PasswordResetExpires   *time.Time `json:"password_reset_expires"`
	StudentIDVerified      bool       `json:"student_id_verified"`
	StudentIDVerification  *string    `json:"student_id_verification"`
	InstitutionName        *string    `json:"institution_name"`
	Gender                 string     `json:"gender"`
	AdditionalNotes        *string    `json:"additional_notes"`
	TwoFAEnabled           bool       `json:"twofa_enabled"`
	TwoFAImage             *string    `json:"twofa_image"`
	TwoFASecret            *string    `json:"twofa_secret"`
}
