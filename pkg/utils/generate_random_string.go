package utils

import "math/rand"

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomOTPCode GenerateRandomOTP generates a random 6-digit OTP
func GenerateRandomOTPCode() string {
	otp := make([]byte, 6) // 6 digits OTP
	_, _ = rand.Read(otp)
	for i := 0; i < len(otp); i++ {
		otp[i] = otp[i]%10 + 48 // Convert to ASCII

		// Ensure that the generated OTP is a 6-digit number
		if otp[i] < 48 || otp[i] > 57 {
			i--
		}
	}

	return string(otp)
}

// GenerateRandomTokenOtp is a secret that is used to generate a random OTP
func GenerateRandomTokenOtp() string {
	// Generate cryptographically secure random bytes
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return ""
	}

	// Filter out non-alphanumeric characters using a mask
	var mask byte = 0b111111 // Mask to keep only alphanumeric characters (lower and uppercase letters, digits)
	for i := range token {
		token[i] &= mask
	}

	// Select random characters from the alphanumeric charset
	for i, b := range token {
		token[i] = alphanumericCharset[b%byte(len(alphanumericCharset))]

		// Ensure that the generated token is 32 characters long
		if i == 31 {
			break
		}
	}

	return string(token)
}
