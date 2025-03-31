package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"github.com/google/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
	"strings"
)

type S3Config struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Minio     bool   `yaml:"minio"`
}

type S3Client struct {
	Client   *s3.Client
	Config   S3Config
}

type Storage interface{
	UploadFile(file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(fileURL string) error
	// GetStorageBucket() string
	// GetStorageEndpoint() string
}

func NewStorage(cfg S3Config) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AccessKey,
				SecretAccessKey: cfg.SecretKey,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации S3: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.Minio
	})

	return &S3Client{
		Client:   client,
		Config:   cfg,
	}, nil
}

func (s *S3Client) UploadFile(file multipart.File, originalFilename string) (string, error) {
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	fileID := uuid.New().String()
	ext := filepath.Ext(originalFilename)
	fileKey := fmt.Sprintf("uploads/%s%s", fileID, ext)

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to s3: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", s.Config.Endpoint, s.Config.Bucket, fileKey)
	return fileURL, nil
}

func (s *S3Client) DeleteFileByURL(fileURL string) error {
	baseURL := fmt.Sprintf("%s/%s/", s.Config.Endpoint, s.Config.Bucket)
	if !strings.HasPrefix(fileURL, baseURL) {
		return fmt.Errorf("wrong URL: %s", fileURL)
	}
	fileKey := strings.TrimPrefix(fileURL, baseURL)

	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return fmt.Errorf("could not delete object from s3: %w", err)
	}
	return nil
}