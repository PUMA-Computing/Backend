package services

import (
	"Backend/pkg/utils"
	"bytes"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type FilesService struct {
}

func NewFilesService() *FilesService {
	return &FilesService{}
}

func (fs *FilesService) IsImageFile(file *multipart.FileHeader) bool {
	openedFile, err := file.Open()
	if err != nil {
		return false
	}
	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			return
		}
	}(openedFile)

	buffer := make([]byte, 512)
	_, err = openedFile.Read(buffer)
	if err != nil {
		return false
	}

	mimeType := http.DetectContentType(buffer)
	return strings.HasPrefix(mimeType, "image/")
}

func (fs *FilesService) IsFileSizeValid(file *multipart.FileHeader) bool {
	return file.Size <= utils.MaxFileSize
}

func (fs *FilesService) IsFileExtensionValid(file *multipart.FileHeader) bool {
	return utils.AllowedImageExtension[strings.ToLower(utils.GetFileExtension(file.Filename))]
}

func (fs *FilesService) IsFileSignatureValid(file *multipart.FileHeader) bool {
	magicNumbers := map[string][]byte{
		".jpg": {0xFF, 0xD8, 0xFF},
		".png": {0x89, 0x50, 0x4E, 0x47},
		".gif": {0x47, 0x49, 0x46, 0x38},
		".bmp": {0x42, 0x4D},
	}

	ext := filepath.Ext(file.Filename)
	magicNumber, ok := magicNumbers[ext]
	if !ok {
		return false
	}

	openedFile, err := file.Open()
	if err != nil {
		return false
	}
	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			return
		}
	}(openedFile)

	buffer := make([]byte, len(magicNumber))
	_, err = openedFile.Read(buffer)
	if err != nil {
		return false
	}

	return bytes.Equal(buffer, magicNumber)
}
