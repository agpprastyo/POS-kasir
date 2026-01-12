package auth

import (
	"POS-kasir/internal/common"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"

	"context"

	"github.com/google/uuid"
)

type AthRepo struct {
	log   logger.ILogger
	minio minio.IMinio
	r2    cloudflarer2.IR2
}

func NewAuthRepo(log logger.ILogger, minio minio.IMinio, r2 cloudflarer2.IR2) IAthRepo {
	return &AthRepo{
		log:   log,
		minio: minio,
		r2:    r2,
	}
}

type IAthRepo interface {
	UploadAvatar(ctx context.Context, filename string, data []byte) (string, error)
	AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error)
}

func (r *AthRepo) UploadAvatar(ctx context.Context, filename string, data []byte) (string, error) {
	// 1. Upload to Minio (Primary)
	minioURL, err := r.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Errorf("UploadAvatar | Failed to upload avatar to Minio: %v", err)
		return "", common.ErrUploadAvatar
	}

	// 2. Upload to R2 (Secondary - Fire and Forget / Log only)
	// We run this synchronously for now to keep it simple, or we could go routine it.
	// Since user wants to log success, we'll do it here.
	r2URL, err := r.r2.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Warnf("UploadAvatar | Failed to upload avatar to R2 (non-blocking): %v", err)
	} else {
		r.log.Infof("UploadAvatar | Successfully uploaded avatar to R2. URL: %s", r2URL)
	}

	// 3. Return Minio URL as requested
	return minioURL, nil
}

func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	r.log.Infof("AvatarLink | AvatarLink called for user %s with avatar %s", userID.String(), avatar)
	if avatar == "" {
		return "", common.ErrAvatarNotFound
	}

	// Priority: Minio
	url, err := r.minio.GetFileShareLink(ctx, avatar)
	if err != nil {
		r.log.Errorf("AvatarLink | Failed to get avatar link from Minio: %v", err)

		// Fallback to R2 if Minio fails? Or strictly stick to Minio?
		// For robustness, let's try R2 if Minio fails, logging the attempt.
		r.log.Warn("AvatarLink | Attempting fallback to R2...")
		r2Url, r2Err := r.r2.GetFileShareLink(ctx, avatar)
		if r2Err != nil {
			r.log.Errorf("AvatarLink | Failed to get avatar link from R2 as well: %v", r2Err)
			return "", common.ErrAvatarLink
		}
		return r2Url, nil
	}
	return url, nil
}
