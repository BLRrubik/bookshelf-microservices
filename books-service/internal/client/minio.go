package client

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const publicPolicy = `
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::bookshelf-covers/*"]
    }
  ]
}
`

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

	return &MinIOClient{
		client:         client,
		bucket:         bucket,
		publicEndpoint: publicEndpoint,
	}, nil
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

	if err = c.client.SetBucketPolicy(ctx, c.bucket, publicPolicy); err != nil {
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

func (c *MinIOClient) UploadFile(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	if size < 0 {
		return fmt.Errorf("object size must be known")
	}

	_, err := c.client.PutObject(
		ctx,
		c.bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fmt.Errorf("put object %q: %w", objectName, err)
	}
	return nil
}

func (c *MinIOClient) GetFileURL(objectName string) string {
	return fmt.Sprintf("%s/%s/%s", c.publicEndpoint, c.bucket, objectName)
}

func (c *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(
		ctx,
		c.bucket,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("remove object %q: %w", objectName, err)
	}

	return nil
}

func GetContentType(filename string) string {
	splitted := strings.Split(filename, ".")

	switch strings.ToLower(splitted[len(splitted)-1]) {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
