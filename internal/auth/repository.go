package auth

import (
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"

	"context"
	"fmt"
	"github.com/google/uuid"
)

type AthRepo struct {
	log   *logger.Logger
	minio minio.IMinio
}

func (r *AthRepo) UploadAvatar(ctx context.Context, filename string, data []byte) (string, error) {
	if r.minio == nil {
		r.log.Error("minio client is nil in UploadAvatar")
		return "", fmt.Errorf("internal error: minio client is not initialized")
	}
	fmt.Println("UploadAvatar repo 1")
	url, err := r.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		fmt.Println("UploadAvatar repo 2")
		r.log.Errorf("Failed to upload avatar: %v", err)
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}
	return url, nil
}
func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *AthRepo) DeleteAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthRepo(log *logger.Logger, minio minio.IMinio) IAthRepo {
	return &AthRepo{
		log:   log,
		minio: minio,
	}
}

type IAthRepo interface {
	UploadAvatar(ctx context.Context, filename string, data []byte) (string, error)
	AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID, avatar string) error
}
