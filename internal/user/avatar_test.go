package user_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAthRepo_UploadAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockR2 := mocks.NewMockIR2(ctrl)

	repo := user.NewAuthRepo(mockLogger, mockR2)

	ctx := context.Background()
	filename := "test.jpg"
	data := []byte("test-data")
	contentType := "image/jpeg"
	expectedURL := "https://example.com/test.jpg"

	t.Run("Success", func(t *testing.T) {
		mockR2.EXPECT().
			UploadFile(ctx, filename, data, contentType).
			Return(expectedURL, nil)

		mockLogger.EXPECT().
			Infof("UploadAvatar | Successfully uploaded avatar to R2. URL: %s", expectedURL)

		url, err := repo.UploadAvatar(ctx, filename, data)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
	})

	t.Run("UploadFailure", func(t *testing.T) {
		uploadErr := errors.New("upload failed")
		mockR2.EXPECT().
			UploadFile(ctx, filename, data, contentType).
			Return("", uploadErr)

		mockLogger.EXPECT().
			Errorf("UploadAvatar | Failed to upload avatar to R2: %v", uploadErr)

		url, err := repo.UploadAvatar(ctx, filename, data)
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.Equal(t, common.ErrUploadAvatar, err)
	})
}

func TestAthRepo_AvatarLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockR2 := mocks.NewMockIR2(ctrl)

	repo := user.NewAuthRepo(mockLogger, mockR2)

	ctx := context.Background()
	userID := uuid.New()
	avatar := "avatar.jpg"
	expectedURL := "https://example.com/avatar.jpg"

	t.Run("Success", func(t *testing.T) {
		mockLogger.EXPECT().
			Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), avatar)

		mockR2.EXPECT().
			GetFileShareLink(ctx, avatar).
			Return(expectedURL, nil)

		url, err := repo.AvatarLink(ctx, userID, avatar)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
	})

	t.Run("EmptyAvatar", func(t *testing.T) {
		mockLogger.EXPECT().
			Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), "")

		url, err := repo.AvatarLink(ctx, userID, "")
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.Equal(t, common.ErrAvatarNotFound, err)
	})

	t.Run("GetLinkFailure", func(t *testing.T) {
		linkErr := errors.New("get link failed")

		mockLogger.EXPECT().
			Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), avatar)

		mockR2.EXPECT().
			GetFileShareLink(ctx, avatar).
			Return("", linkErr)

		mockLogger.EXPECT().
			Errorf("AvatarLink | Failed to get avatar link from R2: %v", linkErr)

		url, err := repo.AvatarLink(ctx, userID, avatar)
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.Equal(t, common.ErrAvatarLink, err)
	})
}
