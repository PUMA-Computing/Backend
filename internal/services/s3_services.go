package services

import (
	"Backend/configs"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client *s3.Client
	bucket   string
}

func NewS3Service() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	var region = s3Config.AWSRegion
	var bucket = s3Config.AWSBucket
	var accessKey = s3Config.AwsAccessKeyId
	var secretKey = s3Config.AwsSecretAccessKey

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client access key and so on
	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (s *S3Service) UploadFile(ctx context.Context, key string, file []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	}

	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// GetBucket file from S3
func (s *S3Service) GetFile(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	output, err := s.s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
