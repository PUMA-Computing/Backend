package utils

import (
	"bytes"
	"encoding/base64"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"log"
)

func GenerateTOTPKey(email string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "PUFA Computing",
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA256,
	})
	if err != nil {
		return nil, err
	}
	log.Println("Generated TOTP key:", key)
	return key, nil
}

func GenerateQRCodeBase64(key *otp.Key) (string, error) {
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	qrCode := base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrCode, nil
}
