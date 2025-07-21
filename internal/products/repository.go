package products

import (
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"
	"context"
)

type PrdRepo struct {
	minio minio.IMinio
	log   logger.ILogger
}

func (p PrdRepo) UploadImageToMinio(ctx context.Context, filename string, data []byte) (string, error) {
	url, err := p.minio.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		p.log.Errorf("Failed to upload product image: %v", err)
		return "", err
	}
	return url, nil
}

func (p PrdRepo) PrdImageLink(ctx context.Context, prdID string, image string) (string, error) {
	p.log.Infof("PrdImageLink called for product %s with image %s", prdID, image)
	if image == "" {
		return "", nil // No image provided, return empty string
	}

	url, err := p.minio.GetFileShareLink(ctx, image)
	if err != nil {
		p.log.Errorf("Failed to get product image link: %v", err)
		return "", err
	}
	return url, nil
}

type IPrdRepo interface {
	UploadImageToMinio(ctx context.Context, filename string, data []byte) (string, error)
	PrdImageLink(ctx context.Context, prdID string, image string) (string, error)
}

func NewPrdRepo(minio minio.IMinio, log logger.ILogger) IPrdRepo {
	return &PrdRepo{
		minio: minio,
		log:   log,
	}
}
