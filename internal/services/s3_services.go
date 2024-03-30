package services

import (
	"Backend/configs"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type S3Service struct {
	s3Client *s3.Client
	bucket   string
}

func NewAWSService() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	var region = s3Config.AWSRegion
	var bucket = s3Config.S3Bucket
	var accessKey = s3Config.AWSAccessKeyId
	var secretKey = s3Config.AWSSecretAccessKey

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

func NewR2Service() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	var bucket = s3Config.S3Bucket
	var accountID = s3Config.CloudflareAccountId
	var accessKey = s3Config.CloudflareR2AccessId
	var secretKey = s3Config.CloudflareR2AccessKey

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client access key and so on
	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (s *S3Service) UploadFileToAWS(ctx context.Context, key string, file []byte) error {
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

func (s *S3Service) UploadFileToR2(ctx context.Context, key string, file []byte) error {
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

// GetFileAWS GetFile GetBucket file from S3
func (s *S3Service) GetFileAWS(ctx context.Context, directory, slug string) (string, error) {

	key := directory + "/" + slug + ".jpg"

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.HeadObject(ctx, input)
	if err != nil {
		return "", err
	}

	return "https://id.pufacomputing.live/" + key, nil
}

func (s *S3Service) GetFileR2(ctx context.Context, directory, slug string) (string, error) {

	key := directory + "/" + slug + ".jpg"

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.HeadObject(ctx, input)
	if err != nil {
		return "", err
	}

	return "https://sg.pufacomputing.live/" + key, nil
}
