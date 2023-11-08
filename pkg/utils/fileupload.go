package utils

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
)

const MaxFileSize = 20 * 1024 * 1024

var AllowedImageExtension = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func GenerateUniqueFileName(originalName string) string {
	randomToken := make([]byte, 16)
	_, err := rand.Read(randomToken)
	if err != nil {
		panic(err)
	}

	uniqueToken := hex.EncodeToString(randomToken)
	fileExtension := filepath.Ext(originalName)
	uniqueFileName := uniqueToken + fileExtension

	return uniqueFileName
}

func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}
