package services

import (
	"Backend/configs"
	"Backend/pkg/utils"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime"
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

func (fs *FilesService) UploadFileToBucket(file *multipart.FileHeader, fileName string) (err error) {
	loadConfig := configs.LoadConfig()
	var accountId = loadConfig.CloudflareAccountId
	var accessKeyId = loadConfig.CloudflareR2AccessId
	var accessKeySecret = loadConfig.CloudflareR2AccessKey
	var url = fmt.Sprintf("https://%s.r2.cloudflarestorage.com/", accountId)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: url,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	contentType := getContentType(fileName)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(loadConfig.CloudflareR2Bucket),
		Key:         aws.String("images/" + fileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}, func(options *s3.Options) {
		options.Region = "APAC"
	})

	if err != nil {
		return err
	}
	return nil
}

func getContentType(filePath string) string {
	// Use the mime package to determine the content type based on the file extension
	ext := mime.TypeByExtension(filepath.Ext(filePath))
	if ext != "" {
		return ext
	}
	// If mime.TypeByExtension fails, you can provide a default content type
	return "application/octet-stream"
}
