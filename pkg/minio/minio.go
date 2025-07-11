package minio

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

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
		log.Error("Failed to create Minio client", "error", err)
		return nil, err
	}

	log.Println("Created Minio client")

	return &Minio{
		Cfg:    cfg,
		Log:    log,
		Client: client,
	}, nil
}

type IMinio interface {
	UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error)
	GetFileShareLink(ctx context.Context, objectName string) (string, error)
}

// GetFileShareLink generates a presigned URL for sharing a file.
func (m *Minio) GetFileShareLink(ctx context.Context, objectName string) (string, error) {

	m.Log.Info("GetFileShareLink", "objectName", objectName)
	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		m.Cfg.Minio.Bucket,
		objectName,
		time.Duration(m.Cfg.Minio.ExpirySec)*time.Second,
		nil,
	)
	if err != nil {
		m.Log.Error("Failed to generate presigned URL", "error", err, "objectName", objectName)
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *Minio) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	if m == nil || m.Client == nil || m.Cfg == nil {
		if m != nil && m.Log != nil {
			m.Log.Error("Minio or its dependencies are nil")
		}
		return "", fmt.Errorf("minio or its dependencies are nil")
	}
	fmt.Println("UploadFile repo 1")
	reader := bytes.NewReader(data)
	fmt.Println("UploadFile repo 2")
	_, err := m.Client.PutObject(ctx, m.Cfg.Minio.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	fmt.Println("UploadFile repo 3")
	if err != nil {
		m.Log.Error("Failed to upload file to Minio", "error", err, "objectName", objectName)
		return "", err
	}

	fmt.Println("UploadFile repo 4")

	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		m.Cfg.Minio.Bucket,
		objectName,
		time.Duration(m.Cfg.Minio.ExpirySec)*time.Second,
		nil,
	)

	fmt.Println("UploadFile repo 5")
	if err != nil {
		m.Log.Error("Failed to generate presigned URL after upload", "error", err, "objectName", objectName)
		return "", err
	}

	m.Log.Info("File uploaded successfully", "objectName", objectName, "url", presignedURL.String())
	return presignedURL.String(), nil
}
