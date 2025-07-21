package auth

import (
	"POS-kasir/mocks"
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/google/uuid"
)

func TestUploadAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockIMinio(ctrl)

	mockMinio.EXPECT().
		UploadFile(gomock.Any(), "avatar.jpg", gomock.Any(), "image/jpeg").
		Return("http://example.com/avatar.jpg", nil)

	repo := &AthRepo{
		log:   &mocks.MockILogger{},
		minio: mockMinio,
	}
	url, err := repo.UploadAvatar(context.Background(), "avatar.jpg", []byte("data"))
	if err != nil || url != "http://example.com/avatar.jpg" {
		t.Errorf("expected success, got err=%v, url=%s", err, url)
	}
}

func TestUploadAvatar_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockIMinio(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)

	// Expect Errorf to be called with any arguments
	mockLogger.EXPECT().
		Errorf(gomock.Any(), gomock.Any()).
		Times(1)

	mockMinio.EXPECT().
		UploadFile(gomock.Any(), "avatar.jpg", gomock.Any(), "image/jpeg").
		Return("", fmt.Errorf("upload error"))

	repo := &AthRepo{
		log:   mockLogger,
		minio: mockMinio,
	}
	url, err := repo.UploadAvatar(context.Background(), "avatar.jpg", []byte("data"))
	if err == nil || url != "" {
		t.Errorf("expected error, got err=%v, url=%s", err, url)
	}
}

func TestAvatarLink_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockIMinio(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	id := uuid.New()

	// Expect Infof to be called (optional, for coverage)
	mockLogger.EXPECT().
		Infof(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	// Expect GetFileShareLink to be called and return success
	mockMinio.EXPECT().
		GetFileShareLink(gomock.Any(), "avatar.jpg").
		Return("http://example.com/avatar.jpg", nil)

	repo := &AthRepo{
		log:   mockLogger,
		minio: mockMinio,
	}
	url, err := repo.AvatarLink(context.Background(), id, "avatar.jpg")
	if err != nil || url != "http://example.com/avatar.jpg" {
		t.Errorf("expected success, got err=%v, url=%s", err, url)
	}
}

func TestAvatarLink_EmptyAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockILogger(ctrl)
	mockMinio := mocks.NewMockIMinio(ctrl)
	id := uuid.New()

	// Expect Infof to be called for logging
	mockLogger.EXPECT().
		Infof(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	repo := &AthRepo{
		log:   mockLogger,
		minio: mockMinio,
	}
	url, err := repo.AvatarLink(context.Background(), id, "")
	if err == nil || url != "" {
		t.Errorf("expected error for empty avatar, got err=%v, url=%s", err, url)
	}
}

func TestAvatarLink_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockILogger(ctrl)
	mockMinio := mocks.NewMockIMinio(ctrl)
	id := uuid.New()

	// Expect Infof to be called for logging
	mockLogger.EXPECT().
		Infof(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	// Expect GetFileShareLink to return an error
	mockMinio.EXPECT().
		GetFileShareLink(gomock.Any(), "avatar.jpg").
		Return("", fmt.Errorf("link error"))

	// Expect Errorf to be called for error logging
	mockLogger.EXPECT().
		Errorf(gomock.Any(), gomock.Any()).
		Times(1)

	repo := &AthRepo{
		log:   mockLogger,
		minio: mockMinio,
	}
	url, err := repo.AvatarLink(context.Background(), id, "avatar.jpg")
	if err == nil || url != "" {
		t.Errorf("expected error, got err=%v, url=%s", err, url)
	}
}
