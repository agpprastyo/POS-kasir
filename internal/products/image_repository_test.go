package products_test

import (
	"POS-kasir/internal/products"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductImageRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockR2 := mocks.NewMockIR2(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	repo := products.NewProductImageRepository(mockR2, mockLogger)
	ctx := context.Background()

	t.Run("UploadImage_Success", func(t *testing.T) {
		mockR2.EXPECT().UploadFile(ctx, "test.jpg", []byte("data"), "image/jpeg").Return("http://r2.com/test.jpg", nil)
		url, err := repo.UploadImage(ctx, "test.jpg", []byte("data"))
		assert.NoError(t, err)
		assert.Equal(t, "http://r2.com/test.jpg", url)
	})

	t.Run("UploadImage_Error", func(t *testing.T) {
		mockR2.EXPECT().UploadFile(ctx, "test.jpg", []byte("data"), "image/jpeg").Return("", errors.New("upload fail"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
		url, err := repo.UploadImage(ctx, "test.jpg", []byte("data"))
		assert.Error(t, err)
		assert.Empty(t, url)
	})

	t.Run("PrdImageLink_Success", func(t *testing.T) {
		mockR2.EXPECT().GetFileShareLink(ctx, "test.jpg").Return("http://signed-url.com", nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		url, err := repo.PrdImageLink(ctx, "123", "test.jpg")
		assert.NoError(t, err)
		assert.Equal(t, "http://signed-url.com", url)
	})

	t.Run("PrdImageLink_EmptyImage", func(t *testing.T) {
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		url, err := repo.PrdImageLink(ctx, "123", "")
		assert.NoError(t, err)
		assert.Empty(t, url)
	})

	t.Run("PrdImageLink_Error", func(t *testing.T) {
		mockR2.EXPECT().GetFileShareLink(ctx, "test.jpg").Return("", errors.New("link fail"))
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
		url, err := repo.PrdImageLink(ctx, "123", "test.jpg")
		assert.Error(t, err)
		assert.Empty(t, url)
	})

	t.Run("R2NotInitialized", func(t *testing.T) {
		repoNoR2 := products.NewProductImageRepository(nil, mockLogger)
		
		t.Run("Upload", func(t *testing.T) {
			mockLogger.EXPECT().Errorf(gomock.Any()).AnyTimes()
			url, err := repoNoR2.UploadImage(ctx, "test.jpg", []byte("data"))
			assert.NoError(t, err)
			assert.Empty(t, url)
		})

		t.Run("Link", func(t *testing.T) {
			mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
			url, err := repoNoR2.PrdImageLink(ctx, "123", "test.jpg")
			assert.NoError(t, err)
			assert.Empty(t, url)
		})
	})
}
