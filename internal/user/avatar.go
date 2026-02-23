package user

import (
	"POS-kasir/internal/common"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/logger"

	"context"

	"github.com/google/uuid"
)

type AthRepo struct {
	log logger.ILogger
	r2  cloudflarer2.IR2
}

func NewAuthRepo(log logger.ILogger, r2 cloudflarer2.IR2) IAthRepo {
	return &AthRepo{
		log: log,
		r2:  r2,
	}
}

type IAthRepo interface {
	UploadAvatar(ctx context.Context, filename string, data []byte) (string, error)
	AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error)
}

func (r *AthRepo) UploadAvatar(ctx context.Context, filename string, data []byte) (string, error) {
	if r.r2 == nil {
		r.log.Errorf("UploadAvatar | R2 storage is not initialized")
		return "", common.ErrUploadAvatar
	}

	// Upload to R2
	r2URL, err := r.r2.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Errorf("UploadAvatar | Failed to upload avatar to R2: %v", err)
		return "", common.ErrUploadAvatar
	}

	r.log.Infof("UploadAvatar | Successfully uploaded avatar to R2. URL: %s", r2URL)

	return r2URL, nil
}

func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	r.log.Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), avatar)
	if avatar == "" {
		return "", common.ErrAvatarNotFound
	}
	if r.r2 == nil {
		r.log.Warnf("AvatarLink | R2 storage is not initialized, returning empty string for avatar %s", avatar)
		return "", nil
	}

	// Get link from R2
	r2Url, err := r.r2.GetFileShareLink(ctx, avatar)
	if err != nil {
		r.log.Errorf("AvatarLink | Failed to get avatar link from R2: %v", err)
		return "", common.ErrAvatarLink
	}
	return r2Url, nil
}
