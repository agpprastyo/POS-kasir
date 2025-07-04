package minio

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"net/url"

	"time"
)

type Minio struct {
	Cfg    *config.AppConfig
	Log    *logger.Logger
	Client *minio.Client
}

func NewMinio(cfg *config.AppConfig, log *logger.Logger) (*Minio, error) {
	client, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Minio{
		Cfg:    cfg,
		Log:    log,
		Client: client,
	}, nil
}

type IMinio interface {
	UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error)
	GetFileShareLink(ctx context.Context, objectName string, expirySeconds int64) (string, error)
}

// GetFileShareLink generates a presigned URL for sharing a file.
func (m *Minio) GetFileShareLink(ctx context.Context, objectName string, expirySeconds int64) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		m.Cfg.Minio.Bucket,
		objectName,
		time.Duration(expirySeconds)*time.Second,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// In pkg/minio/minio.go (or similar)
func (m *Minio) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := m.Client.PutObject(ctx, m.Cfg.Minio.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		m.Cfg.Minio.Bucket,
		objectName,
		time.Hour,
		nil,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
