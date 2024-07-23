package services

import (
	"Backend/pkg/utils"
	"context"
	"github.com/google/uuid"
	"log"
	"time"
)

type OTPManager struct {
}

func NewOTPManager() *OTPManager {
	return &OTPManager{}
}

func (om *OTPManager) GenerateOTP(userID uuid.UUID, token string, expiration time.Duration) (string, error) {
	otpCode := utils.GenerateRandomOTPCode()
	expiresAt := time.Now().Add(expiration)

	ctx := context.Background()
	otpKey := "otp:" + userID.String()

	err := utils.Rdb.HSet(ctx, otpKey, map[string]interface{}{
		"TokenOTP":  token,
		"OTPCode":   otpCode,
		"ExpiresAt": expiresAt.Format(time.RFC3339),
	}).Err()
	if err != nil {
		return "", err
	}

	err = utils.Rdb.Expire(ctx, otpKey, expiration).Err()
	if err != nil {
		return "", err
	}

	return otpCode, nil
}

func (om *OTPManager) VerifyOTP(userID uuid.UUID, tokenOtp, otpCode string) bool {
	ctx := context.Background()
	otpKey := "otp:" + userID.String()

	otpData, err := utils.Rdb.HGetAll(ctx, otpKey).Result()
	if err != nil {
		log.Printf("Error retrieving OTP data: %v", err)
		return false
	}
	if len(otpData) == 0 {
		return false
	}

	storedTokenOTP := otpData["TokenOTP"]
	storedOTPCode := otpData["OTPCode"]
	expiresAt, err := time.Parse(time.RFC3339, otpData["ExpiresAt"])
	if err != nil {
		log.Printf("Error parsing OTP expiration time: %v", err)
		return false
	}

	if storedTokenOTP != tokenOtp {
		return false
	}

	if time.Now().After(expiresAt) {
		return false
	}

	return storedOTPCode == otpCode
}
