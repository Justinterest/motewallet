package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"motewallet/internal/config"
)

type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	prefix        string
	enabled       bool
}

func NewS3Storage(cfg config.S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return &S3Storage{enabled: false}, nil
	}

	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	var client *s3.Client
	if cfg.Endpoint != "" {
		client = s3.New(s3.Options{
			Region:           cfg.Region,
			Credentials:      awsCfg.Credentials,
			UsePathStyle:     cfg.ForcePathStyle,
			BaseEndpoint:     aws.String(cfg.Endpoint),
		})
	} else {
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.UsePathStyle = cfg.ForcePathStyle
		})
	}

	prefix := strings.Trim(cfg.Prefix, "/")
	if prefix != "" {
		prefix += "/"
	}

	return &S3Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
		prefix:        prefix,
		enabled:       true,
	}, nil
}

func (s *S3Storage) Enabled() bool {
	return s.enabled
}

func (s *S3Storage) ObjectKey(merchantID uint64, filename string) string {
	ext := fileExtension(filename)
	return fmt.Sprintf("%skyc/%d/%s%s", s.prefix, merchantID, uuid.NewString(), ext)
}

func (s *S3Storage) ValidateObjectKey(merchantID uint64, objectKey string) bool {
	expectedPrefix := fmt.Sprintf("%skyc/%d/", s.prefix, merchantID)
	return strings.HasPrefix(objectKey, expectedPrefix)
}

func (s *S3Storage) PresignPut(ctx context.Context, objectKey, contentType string, expires time.Duration) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("S3 storage is not configured")
	}

	out, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return out.URL, nil
}

func (s *S3Storage) PresignGet(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("S3 storage is not configured")
	}

	out, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}
	return out.URL, nil
}

func (s *S3Storage) GetObject(ctx context.Context, objectKey string) ([]byte, string, error) {
	if !s.enabled {
		return nil, "", fmt.Errorf("S3 storage is not configured")
	}

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, "", fmt.Errorf("get object: %w", err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read object body: %w", err)
	}

	contentType := "application/octet-stream"
	if out.ContentType != nil && *out.ContentType != "" {
		contentType = *out.ContentType
	}
	return data, contentType, nil
}

func fileExtension(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if idx := strings.LastIndex(name, "."); idx >= 0 && idx < len(name)-1 {
		ext := strings.ToLower(name[idx:])
		if len(ext) <= 8 {
			return ext
		}
	}
	return ""
}
