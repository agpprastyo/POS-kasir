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
	url, err := r.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Errorf("Failed to upload avatar: %v", err)
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}
	return url, nil
}
func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	r.log.Infof("AvatarLink called for user %s with avatar %s", userID.String(), avatar)
	if avatar == "" {
		return "", fmt.Errorf("avatar is empty for user %s", userID.String())
	}

	url, err := r.minio.GetFileShareLink(ctx, avatar)
	if err != nil {
		r.log.Errorf("Failed to get avatar link: %v", err)
		return "", fmt.Errorf("failed to get avatar link: %w", err)
	}
	return url, nil
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
