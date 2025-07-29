package auth

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"

	"context"

	"github.com/google/uuid"
)

type AthRepo struct {
	log   logger.ILogger
	minio minio.IMinio
}

func NewAuthRepo(log logger.ILogger, minio minio.IMinio) IAthRepo {
	return &AthRepo{
		log:   log,
		minio: minio,
	}
}

type IAthRepo interface {
	UploadAvatar(ctx context.Context, filename string, data []byte) (string, error)
	AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error)
}

func (r *AthRepo) UploadAvatar(ctx context.Context, filename string, data []byte) (string, error) {
	url, err := r.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Errorf("UploadAvatar | Failed to upload avatar: %v", err)
		return "", common.ErrUploadAvatar
	}
	return url, nil
}
func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	r.log.Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), avatar)
	if avatar == "" {
		return "", common.ErrAvatarNotFound
	}

	url, err := r.minio.GetFileShareLink(ctx, avatar)
	if err != nil {
		r.log.Errorf("AvatarLink | Failed to get avatar link: %v", err)
		return "", common.ErrAvatarLink
	}
	return url, nil
}
