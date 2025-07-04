package auth

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"
	"context"
	"github.com/google/uuid"
)

type AthRepo struct {
	repo  repository.Queries
	log   *logger.Logger
	minio minio.IMinio
}

func (r *AthRepo) UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error) {
	filename := "avatars/" + userID.String() + ".jpg"

	url, err := r.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		r.log.Errorf("failed to upload avatar: %v", err)
		return nil, err
	}

	params := repository.UpdateAvatarParams{
		ID:     userID,
		Avatar: &filename,
	}

	err = r.repo.UpdateAvatar(ctx, params)
	if err != nil {
		r.log.Errorf("failed to update user avatar in db: %v", err)
		return nil, err
	}

	profile, err := r.repo.GetUserByID(ctx, userID)
	if err != nil {
		r.log.Errorf("failed to get user profile: %v", err)
		return nil, err
	}

	response := &ProfileResponse{
		Username:  profile.Username,
		Avatar:    &url,
		CreatedAt: profile.CreatedAt.Time,
		UpdatedAt: profile.UpdatedAt.Time,
		Role:      profile.Role,
	}

	return response, nil
}

func (r *AthRepo) AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *AthRepo) DeleteAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	//TODO implement me
	panic("implement me")
}

func (r *AthRepo) NewAuthRepo(repo repository.Queries, log *logger.Logger, minio minio.IMinio) IAthRepo {
	return &AthRepo{
		repo:  repo,
		log:   log,
		minio: minio,
	}
}

type IAthRepo interface {
	UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error)
	AvatarLink(ctx context.Context, userID uuid.UUID, avatar string) (string, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID, avatar string) error
}
