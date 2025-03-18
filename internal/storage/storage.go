package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
)

type S3Client struct {
	Client   *s3.Client
	Config   S3Config
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

func (s *S3Client) UploadFile(file multipart.File, fileKey string) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла в S3: %w", err)
	}
	return nil
}
