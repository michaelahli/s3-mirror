package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/michaelahli/s3-mirror/internal/config"
)

// Client is the interface for storage operations
type Client interface {
	ListObjects(ctx context.Context, prefix string) ([]ObjectInfo, error)
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	PutObject(ctx context.Context, key string, body io.Reader, size int64) error
	HeadObject(ctx context.Context, key string) (*ObjectInfo, error)
	CopyObject(ctx context.Context, sourceKey, targetKey string) error
	GetBucket() string
}

// ObjectInfo contains metadata about an object
type ObjectInfo struct {
	Key          string
	Size         int64
	ETag         string
	LastModified string
}

// NewClient creates a new storage client based on the configuration
func NewClient(cfg config.StorageConfig) (Client, error) {
	switch cfg.Type {
	case "s3":
		return NewS3Client(cfg)
	case "minio":
		return NewMinIOClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
