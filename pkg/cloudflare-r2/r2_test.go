package cloudflarer2_test

import (
	"POS-kasir/config"
	"POS-kasir/mocks"
	r2 "POS-kasir/pkg/cloudflare-r2"
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCloudflareR2_BucketExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockStorageClient(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	cfg := &config.AppConfig{
		CloudflareR2: struct {
			AccountID    string
			AccessKey    string
			SecretKey    string
			Bucket       string
			PublicDomain string
			ExpirySec    int64
		}{
			Bucket: "test-bucket",
		},
	}

	r2Instance := &r2.CloudflareR2{
		Cfg:    cfg,
		Log:    mockLogger,
		Client: mockClient,
	}

	ctx := context.Background()

	t.Run("Success_Exists", func(t *testing.T) {
		mockClient.EXPECT().BucketExists(ctx, "test-bucket").Return(true, nil)
		exists, err := r2Instance.BucketExists(ctx)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Success_NotExists", func(t *testing.T) {
		mockClient.EXPECT().BucketExists(ctx, "test-bucket").Return(false, nil)
		exists, err := r2Instance.BucketExists(ctx)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Error_ClientFailure", func(t *testing.T) {
		expectedErr := errors.New("bucket check failed")
		mockClient.EXPECT().BucketExists(ctx, "test-bucket").Return(false, expectedErr)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		exists, err := r2Instance.BucketExists(ctx)
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Error_NilClient", func(t *testing.T) {
		nilClientR2 := &r2.CloudflareR2{
			Cfg:    cfg,
			Log:    mockLogger,
			Client: nil,
		}
		exists, err := nilClientR2.BucketExists(ctx)
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "r2 client or dependencies are nil")
	})
}

func TestCloudflareR2_GetFileShareLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockStorageClient(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	cfg := &config.AppConfig{
		CloudflareR2: struct {
			AccountID    string
			AccessKey    string
			SecretKey    string
			Bucket       string
			PublicDomain string
			ExpirySec    int64
		}{
			Bucket:    "test-bucket",
			ExpirySec: 3600,
		},
	}

	r2Instance := &r2.CloudflareR2{
		Cfg:    cfg,
		Log:    mockLogger,
		Client: mockClient,
	}

	ctx := context.Background()
	objectName := "image.jpg"

	t.Run("Success_PublicDomain", func(t *testing.T) {
		// Create a separate config for this test to avoid modifying the shared one safely
		publicCfg := *cfg
		publicCfg.CloudflareR2.PublicDomain = "https://pub.example.com"

		publicR2 := &r2.CloudflareR2{
			Cfg:    &publicCfg,
			Log:    mockLogger,
			Client: mockClient,
		}

		url, err := publicR2.GetFileShareLink(ctx, objectName)
		assert.NoError(t, err)
		assert.Equal(t, "https://pub.example.com/image.jpg", url)
	})

	t.Run("Success_PresignedURL", func(t *testing.T) {
		// Ensure PublicDomain is empty
		cfg.CloudflareR2.PublicDomain = ""

		expectedURL, _ := url.Parse("https://r2.example.com/image.jpg?token=123")

		mockClient.EXPECT().
			PresignedGetObject(ctx, "test-bucket", objectName, time.Duration(3600)*time.Second, nil).
			Return(expectedURL, nil)

		urlStr, err := r2Instance.GetFileShareLink(ctx, objectName)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL.String(), urlStr)
	})

	t.Run("Error_PresignedURLFailure", func(t *testing.T) {
		cfg.CloudflareR2.PublicDomain = ""
		expectedErr := errors.New("presign failed")

		mockClient.EXPECT().
			PresignedGetObject(ctx, "test-bucket", objectName, time.Duration(3600)*time.Second, nil).
			Return(nil, expectedErr)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		urlStr, err := r2Instance.GetFileShareLink(ctx, objectName)
		assert.Error(t, err)
		assert.Empty(t, urlStr)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCloudflareR2_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockStorageClient(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	cfg := &config.AppConfig{
		CloudflareR2: struct {
			AccountID    string
			AccessKey    string
			SecretKey    string
			Bucket       string
			PublicDomain string
			ExpirySec    int64
		}{
			Bucket: "test-bucket",
		},
	}

	r2Instance := &r2.CloudflareR2{
		Cfg:    cfg,
		Log:    mockLogger,
		Client: mockClient,
	}

	ctx := context.Background()
	objectName := "upload.jpg"
	data := []byte("file content")
	contentType := "image/jpeg"

	t.Run("Success", func(t *testing.T) {
		// Expect PutObject
		mockClient.EXPECT().
			PutObject(ctx, "test-bucket", objectName, gomock.Any(), int64(len(data)), minio.PutObjectOptions{ContentType: contentType}).
			Return(minio.UploadInfo{}, nil)

		// Expect GetFileShareLink behavior (assuming public domain is empty, it calls PresignedGetObject)
		expectedURL, _ := url.Parse("https://r2.example.com/upload.jpg")
		mockClient.EXPECT().
			PresignedGetObject(ctx, "test-bucket", objectName, gomock.Any(), nil).
			Return(expectedURL, nil)

		urlStr, err := r2Instance.UploadFile(ctx, objectName, data, contentType)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL.String(), urlStr)
	})

	t.Run("Error_NilClient", func(t *testing.T) {
		nilClientR2 := &r2.CloudflareR2{
			Cfg:    cfg,
			Log:    mockLogger,
			Client: nil,
		}
		urlStr, err := nilClientR2.UploadFile(ctx, objectName, data, contentType)
		assert.Error(t, err)
		assert.Empty(t, urlStr)
		assert.Contains(t, err.Error(), "r2 client or dependencies are nil")
	})

	t.Run("Error_UploadFailure", func(t *testing.T) {
		expectedErr := errors.New("upload failed")
		mockClient.EXPECT().
			PutObject(ctx, "test-bucket", objectName, gomock.Any(), int64(len(data)), minio.PutObjectOptions{ContentType: contentType}).
			Return(minio.UploadInfo{}, expectedErr)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		urlStr, err := r2Instance.UploadFile(ctx, objectName, data, contentType)
		assert.Error(t, err)
		assert.Empty(t, urlStr)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Error_GetShareLinkFailure", func(t *testing.T) {
		// Expect PutObject Success
		mockClient.EXPECT().
			PutObject(ctx, "test-bucket", objectName, gomock.Any(), int64(len(data)), minio.PutObjectOptions{ContentType: contentType}).
			Return(minio.UploadInfo{}, nil)

		// Expect GetFileShareLink Failure
		expectedErr := errors.New("sign failed")
		mockClient.EXPECT().
			PresignedGetObject(ctx, "test-bucket", objectName, gomock.Any(), nil).
			Return(nil, expectedErr)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		urlStr, err := r2Instance.UploadFile(ctx, objectName, data, contentType)
		assert.Error(t, err)
		assert.Empty(t, urlStr)
		assert.Equal(t, expectedErr, err)
	})
}

func TestNewCloudflareR2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockStorageClient(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	cfg := &config.AppConfig{
		CloudflareR2: struct {
			AccountID    string
			AccessKey    string
			SecretKey    string
			Bucket       string
			PublicDomain string
			ExpirySec    int64
		}{
			AccountID: "test-account",
			AccessKey: "test-key",
			SecretKey: "test-secret",
			Bucket:    "test-bucket",
		},
	}

	// Backup original newClient function
	origNewMinioClient := r2.NewMinioClient
	defer func() { r2.NewMinioClient = origNewMinioClient }()

	t.Run("Success_BucketExists", func(t *testing.T) {
		r2.NewMinioClient = func(endpoint string, opts *minio.Options) (r2.StorageClient, error) {
			return mockClient, nil
		}

		mockClient.EXPECT().BucketExists(gomock.Any(), "test-bucket").Return(true, nil)
		mockLogger.EXPECT().Println("Created Cloudflare R2 client").Times(1)

		instance, err := r2.NewCloudflareR2(cfg, mockLogger)
		assert.NoError(t, err)
		assert.NotNil(t, instance)
	})

	t.Run("Success_BucketNotExists", func(t *testing.T) {
		r2.NewMinioClient = func(endpoint string, opts *minio.Options) (r2.StorageClient, error) {
			return mockClient, nil
		}

		mockClient.EXPECT().BucketExists(gomock.Any(), "test-bucket").Return(false, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)
		mockLogger.EXPECT().Println("Created Cloudflare R2 client").Times(1)

		instance, err := r2.NewCloudflareR2(cfg, mockLogger)
		assert.NoError(t, err)
		assert.NotNil(t, instance)
	})

	t.Run("Error_NewClientFailure", func(t *testing.T) {
		expectedErr := errors.New("init failed")
		r2.NewMinioClient = func(endpoint string, opts *minio.Options) (r2.StorageClient, error) {
			return nil, expectedErr
		}

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		instance, err := r2.NewCloudflareR2(cfg, mockLogger)
		assert.Error(t, err)
		assert.Nil(t, instance)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Error_BucketCheckFailure", func(t *testing.T) {
		expectedErr := errors.New("bucket check failed")
		r2.NewMinioClient = func(endpoint string, opts *minio.Options) (r2.StorageClient, error) {
			return mockClient, nil
		}

		mockClient.EXPECT().BucketExists(gomock.Any(), "test-bucket").Return(false, expectedErr)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		instance, err := r2.NewCloudflareR2(cfg, mockLogger)
		assert.Error(t, err)
		assert.Nil(t, instance)
		assert.Equal(t, expectedErr, err)
	})
}
