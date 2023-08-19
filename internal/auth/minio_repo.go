package auth

import (
	"context"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error)
	GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucket string, fileName string) error
}
