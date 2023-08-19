package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
)

func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Minio.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
