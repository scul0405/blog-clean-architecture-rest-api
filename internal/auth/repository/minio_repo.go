package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type authMinioRepo struct {
	client *minio.Client
}

func NewAuthMinioRepository(client *minio.Client) auth.MinioRepository {
	return &authMinioRepo{client: client}
}

// PutObject upload file to Minio
func (r *authMinioRepo) PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authMinioRepo.PutObject")
	defer span.Finish()

	options := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := r.client.PutObject(ctx, input.BucketName, r.generateFileName(input.Name), input.File, input.Size, options)
	if err != nil {
		return nil, errors.Wrap(err, "authMinioRepo.FileUpload.PutObject")
	}

	return &uploadInfo, err
}

// GetObject download file from Minio
func (r *authMinioRepo) GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authMinioRepo.GetObject")
	defer span.Finish()

	object, err := r.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "authMinioRepo.FileDownload.GetObject")
	}
	return object, nil
}

// RemoveObject delete file from Minio
func (r *authMinioRepo) RemoveObject(ctx context.Context, bucket string, fileName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authMinioRepo.RemoveObject")
	defer span.Finish()

	if err := r.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{}); err != nil {
		return errors.Wrap(err, "authMinioRepo.RemoveObject")
	}
	return nil
}

func (r *authMinioRepo) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
