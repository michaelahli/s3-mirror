package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelahli/s3-mirror/internal/config"
)

// S3Client wraps the AWS S3 client
type S3Client struct {
	s3Client *s3.Client
	bucket   string
	region   string
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg config.StorageConfig) (*S3Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		s3Client: s3.NewFromConfig(awsCfg),
		bucket:   cfg.Bucket,
		region:   cfg.Region,
	}, nil
}

// ListObjects lists all objects in the bucket with an optional prefix
func (c *S3Client) ListObjects(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			objects = append(objects, ObjectInfo{
				Key:          *obj.Key,
				Size:         *obj.Size,
				ETag:         *obj.ETag,
				LastModified: obj.LastModified.String(),
			})
		}
	}

	return objects, nil
}

// GetObject retrieves an object from S3
func (c *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

// PutObject uploads an object to S3
func (c *S3Client) PutObject(ctx context.Context, key string, body io.Reader, size int64) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

// HeadObject gets metadata about an object without downloading it
func (c *S3Client) HeadObject(ctx context.Context, key string) (*ObjectInfo, error) {
	result, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &ObjectInfo{
		Key:          key,
		Size:         *result.ContentLength,
		ETag:         *result.ETag,
		LastModified: result.LastModified.String(),
	}, nil
}

// CopyObject copies an object within S3 (not used for cross-bucket in this implementation)
func (c *S3Client) CopyObject(ctx context.Context, sourceKey, targetKey string) error {
	copySource := c.bucket + "/" + sourceKey
	_, err := c.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(c.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(targetKey),
	})
	return err
}

// GetBucket returns the bucket name
func (c *S3Client) GetBucket() string {
	return c.bucket
}
