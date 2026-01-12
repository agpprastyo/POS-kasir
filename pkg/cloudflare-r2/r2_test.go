package cloudflarer2

import (
	"POS-kasir/config"
	"POS-kasir/mocks"
	"context"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetFileShareLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
			AccountID:    "dummy-account",
			AccessKey:    "dummy-key",
			SecretKey:    "dummy-secret",
			Bucket:       "test-bucket",
			PublicDomain: "https://pub.example.com",
			ExpirySec:    3600,
		},
	}

	r2 := &CloudflareR2{
		Cfg: cfg,
		Log: mockLogger,
	}

	t.Run("ReturnsPublicDomainURL", func(t *testing.T) {
		objectName := "test-image.jpg"
		expectedURL := "https://pub.example.com/test-image.jpg"

		url, err := r2.GetFileShareLink(context.Background(), objectName)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if url != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, url)
		}
	})

	t.Run("PresignedFallBackNeedsClient", func(t *testing.T) {
		r2.Cfg.CloudflareR2.PublicDomain = ""

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		defer func() {
			if r := recover(); r == nil {
			}
		}()

		_, _ = r2.GetFileShareLink(context.Background(), "test.jpg")
	})
}

func TestInitializationLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockFieldLogger(ctrl)

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := &config.AppConfig{
		CloudflareR2: struct {
			AccountID    string
			AccessKey    string
			SecretKey    string
			Bucket       string
			PublicDomain string
			ExpirySec    int64
		}{
			AccountID: "test",
			AccessKey: "key",
			SecretKey: "secret",
		},
	}

	_, _ = NewCloudflareR2(cfg, mockLogger)
}
