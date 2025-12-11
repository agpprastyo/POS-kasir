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
	Log    logger.ILogger
	Client *minio.Client
}

func NewMinio(cfg *config.AppConfig, log logger.ILogger) (IMinio, error) {
	client, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})

	if err != nil {
		log.Errorf("Failed to create Minio client: %v", err)
		return nil, err
	}

	exists, err := client.BucketExists(context.Background(), cfg.Minio.Bucket)
	if err != nil {
		log.Errorf("Failed to check if bucket exists: %v", err)
		return nil, err
	}
	if !exists {
		log.Infof("Bucket %s does not exist, creating...", cfg.Minio.Bucket)
		if err := client.MakeBucket(context.Background(), cfg.Minio.Bucket, minio.MakeBucketOptions{}); err != nil {
			log.Errorf("Failed to create bucket %s: %v", cfg.Minio.Bucket, err)
			return nil, err
		}
		log.Infof("Created bucket: %s", cfg.Minio.Bucket)
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
	BucketExists(ctx context.Context) (bool, error)
}

// BucketExists checks if the configured Minio bucket exists.
func (m *Minio) BucketExists(ctx context.Context) (bool, error) {
	if m == nil || m.Client == nil || m.Cfg == nil {
		if m != nil && m.Log != nil {
			m.Log.Errorf("Minio or its dependencies are nil")
		}
		return false, fmt.Errorf("minio or its dependencies are nil")
	}

	exists, err := m.Client.BucketExists(ctx, m.Cfg.Minio.Bucket)
	if err != nil {
		m.Log.Errorf("Failed to check if bucket exists: %v", err)
		return false, err
	}

	m.Log.Infof("Bucket exists check: bucket=%s, exists=%v", m.Cfg.Minio.Bucket, exists)
	return exists, nil
}

// GetFileShareLink generates a presigned URL for accessing a file stored in the configured Minio bucket.
func (m *Minio) GetFileShareLink(ctx context.Context, objectName string) (string, error) {

	m.Log.Infof("GetFileShareLink: objectName=%s", objectName)
	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		m.Cfg.Minio.Bucket,
		objectName,
		time.Duration(m.Cfg.Minio.ExpirySec)*time.Second,
		nil,
	)
	if err != nil {
		m.Log.Errorf("Failed to generate presigned URL: %v, objectName=%s", err, objectName)
		return "", err
	}
	return presignedURL.String(), nil
}

// UploadFile uploads a file to the configured Minio bucket and returns a presigned URL for accessing the uploaded file.
func (m *Minio) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	if m == nil || m.Client == nil || m.Cfg == nil {
		if m != nil && m.Log != nil {
			m.Log.Errorf("Minio or its dependencies are nil")
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
		m.Log.Errorf("Failed to upload file to Minio: %v, objectName=%s", err, objectName)
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
		m.Log.Errorf("Failed to generate presigned URL after upload: %v, objectName=%s", err, objectName)
		return "", err
	}

	m.Log.Infof("File uploaded successfully: objectName=%s, url=%s", objectName, presignedURL.String())
	return presignedURL.String(), nil
}
