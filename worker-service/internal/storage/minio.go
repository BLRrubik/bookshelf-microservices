package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucket, publicEndpoint string, useSSL bool) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create MinIO client: %w", err)
	}

	return &MinIOStorage{
		client:         client,
		bucket:         bucket,
		publicEndpoint: publicEndpoint,
	}, nil
}

func (s *MinIOStorage) HealthCheck(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", s.bucket, err)
	}

	if !exists {
		return fmt.Errorf("bucket %q does not exist", s.bucket)
	}

	return nil
}

func (s *MinIOStorage) GetFile(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := s.client.GetObject(
		ctx,
		s.bucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("get object %q: %w", objectName, err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object %q: %w", objectName, err)
	}

	return data, nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) error {
	_, err := s.client.PutObject(
		ctx,
		s.bucket,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("upload object %q: %w", objectName, err)
	}

	return nil
}

func (s *MinIOStorage) GetFileURL(objectName string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicEndpoint, s.bucket, objectName)
}
