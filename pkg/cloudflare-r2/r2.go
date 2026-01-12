package cloudflarer2

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type CloudflareR2 struct {
	Cfg    *config.AppConfig
	Log    logger.ILogger
	Client *minio.Client
}

type IR2 interface {
	UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error)
	GetFileShareLink(ctx context.Context, objectName string) (string, error)
	BucketExists(ctx context.Context) (bool, error)
}

func NewCloudflareR2(cfg *config.AppConfig, log logger.ILogger) (IR2, error) {
	// R2 Endpoint: https://<accountid>.r2.cloudflarestorage.com
	endpoint := fmt.Sprintf("%s.r2.cloudflarestorage.com", cfg.CloudflareR2.AccountID)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.CloudflareR2.AccessKey, cfg.CloudflareR2.SecretKey, ""),
		Secure: true,
	})

	if err != nil {
		log.Errorf("Failed to create R2 client: %v", err)
		return nil, err
	}

	// Verify bucket exists
	exists, err := client.BucketExists(context.Background(), cfg.CloudflareR2.Bucket)
	if err != nil {
		log.Errorf("Failed to check if R2 bucket exists: %v", err)
		return nil, err
	}
	if !exists {
		log.Warnf("R2 Bucket %s does not exist. (Note: Automatic creation might not be supported or permitted)", cfg.CloudflareR2.Bucket)
		// Usually we don't auto-create R2 buckets from app, but mimicking minio behavior if needed.
		// For now, just logging warning.
	}

	log.Println("Created Cloudflare R2 client")

	return &CloudflareR2{
		Cfg:    cfg,
		Log:    log,
		Client: client,
	}, nil
}

func (r *CloudflareR2) BucketExists(ctx context.Context) (bool, error) {
	if r == nil || r.Client == nil || r.Cfg == nil {
		return false, fmt.Errorf("r2 client or dependencies are nil")
	}

	exists, err := r.Client.BucketExists(ctx, r.Cfg.CloudflareR2.Bucket)
	if err != nil {
		r.Log.Errorf("Failed to check if R2 bucket exists: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *CloudflareR2) GetFileShareLink(ctx context.Context, objectName string) (string, error) {
	// If PublicDomain is set, return the public URL directly
	if r.Cfg.CloudflareR2.PublicDomain != "" {
		return fmt.Sprintf("%s/%s", r.Cfg.CloudflareR2.PublicDomain, objectName), nil
	}

	// Falls back to presigned URL if no public domain is configured
	presignedURL, err := r.Client.PresignedGetObject(
		ctx,
		r.Cfg.CloudflareR2.Bucket,
		objectName,
		time.Duration(r.Cfg.CloudflareR2.ExpirySec)*time.Second,
		nil,
	)
	if err != nil {
		r.Log.Errorf("Failed to generate R2 presigned URL: %v, objectName=%s", err, objectName)
		return "", err
	}
	return presignedURL.String(), nil
}

func (r *CloudflareR2) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	if r == nil || r.Client == nil || r.Cfg == nil {
		return "", fmt.Errorf("r2 client or dependencies are nil")
	}

	reader := bytes.NewReader(data)
	_, err := r.Client.PutObject(ctx, r.Cfg.CloudflareR2.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		r.Log.Errorf("Failed to upload file to R2: %v, objectName=%s", err, objectName)
		return "", err
	}

	return r.GetFileShareLink(ctx, objectName)
}
