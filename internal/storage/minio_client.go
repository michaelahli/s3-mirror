package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/michaelahli/s3-mirror/internal/config"
)

// MinIOClient wraps the MinIO client
type MinIOClient struct {
	client *minio.Client
	bucket string
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg config.StorageConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// ListObjects lists all objects in the bucket with an optional prefix
func (c *MinIOClient) ListObjects(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	for object := range c.client.ListObjects(ctx, c.bucket, opts) {
		if object.Err != nil {
			return nil, object.Err
		}

		objects = append(objects, ObjectInfo{
			Key:          object.Key,
			Size:         object.Size,
			ETag:         object.ETag,
			LastModified: object.LastModified.String(),
		})
	}

	return objects, nil
}

// GetObject retrieves an object from MinIO
func (c *MinIOClient) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// PutObject uploads an object to MinIO
func (c *MinIOClient) PutObject(ctx context.Context, key string, body io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, c.bucket, key, body, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

// HeadObject gets metadata about an object without downloading it
func (c *MinIOClient) HeadObject(ctx context.Context, key string) (*ObjectInfo, error) {
	stat, err := c.client.StatObject(ctx, c.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &ObjectInfo{
		Key:          stat.Key,
		Size:         stat.Size,
		ETag:         stat.ETag,
		LastModified: stat.LastModified.String(),
	}, nil
}

// CopyObject copies an object within MinIO (not used for cross-bucket in this implementation)
func (c *MinIOClient) CopyObject(ctx context.Context, sourceKey, targetKey string) error {
	src := minio.CopySrcOptions{
		Bucket: c.bucket,
		Object: sourceKey,
	}
	dst := minio.CopyDestOptions{
		Bucket: c.bucket,
		Object: targetKey,
	}

	_, err := c.client.CopyObject(ctx, dst, src)
	return err
}

// GetBucket returns the bucket name
func (c *MinIOClient) GetBucket() string {
	return c.bucket
}
