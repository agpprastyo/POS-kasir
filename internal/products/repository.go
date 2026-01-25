package products

import (
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/logger"
	"context"
)

type PrdRepo struct {
	r2  cloudflarer2.IR2
	log logger.ILogger
}

func (p PrdRepo) UploadImage(ctx context.Context, filename string, data []byte) (string, error) {
	url, err := p.r2.UploadFile(ctx, filename, data, "image/jpeg")
	if err != nil {
		p.log.Errorf("Failed to upload product image to R2: %v", err)
		return "", err
	}
	return url, nil
}

func (p PrdRepo) PrdImageLink(ctx context.Context, prdID string, image string) (string, error) {
	p.log.Infof("PrdImageLink called for product %s with image %s", prdID, image)
	if image == "" {
		return "", nil // No image provided, return empty string
	}

	url, err := p.r2.GetFileShareLink(ctx, image)
	if err != nil {
		p.log.Errorf("Failed to get product image link from R2: %v", err)
		return "", err
	}
	return url, nil
}

type IPrdRepo interface {
	UploadImage(ctx context.Context, filename string, data []byte) (string, error)
	PrdImageLink(ctx context.Context, prdID string, image string) (string, error)
}

func NewPrdRepo(r2 cloudflarer2.IR2, log logger.ILogger) IPrdRepo {
	return &PrdRepo{
		r2:  r2,
		log: log,
	}
}
