package client

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucket, publicEndpoint string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create MinIO client: %w", err)
	}

	return &MinIOClient{client: client, bucket: bucket}, nil
}

func (c *MinIOClient) EnsureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", c.bucket, err)
	}
	if exists {
		return nil
	}

	if err = c.client.MakeBucket(
		ctx,
		c.bucket,
		minio.MakeBucketOptions{},
	); err != nil {
		return fmt.Errorf("create bucket %q: %w", c.bucket, err)
	}

	if err = c.client.SetBucketPolicy(ctx, c.bucket, "public"); err != nil {
		return fmt.Errorf("set bucket policy: %w", err)
	}

	return nil
}

func (c *MinIOClient) HealthCheck(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", c.bucket, err)
	}

	if !exists {
		return fmt.Errorf("bucket %q does not exist", c.bucket)
	}

	return nil
}
