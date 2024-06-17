package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/otel/trace"
)

const (
	DigitalOceanRegion = `us-east-1`
)

type S3Storage struct {
	bucket string
	prefix string
	s3     *s3.Client

	tracer trace.Tracer
}

type S3Options struct {
	Endpoint   string `json:"endpoint"`
	Region     string `json:"region"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
	Bucket     string `json:"bucket"`
	PathPrefix string `json:"path_prefix"`
}

func (opts *S3Options) Validate() error {
	if len(opts.Region) == 0 {
		return fmt.Errorf("region is required")
	}

	if len(opts.Key) == 0 || len(opts.Secret) == 0 {
		return fmt.Errorf("key or secret is required")
	}

	if len(opts.Bucket) == 0 {
		return fmt.Errorf("bucket is required")
	}

	return nil
}

func NewS3Storage(opts *S3Options) (*S3Storage, error) {
	s3Config := aws.Config{
		Credentials:  credentials.NewStaticCredentialsProvider(opts.Key, opts.Secret, ""),
		BaseEndpoint: aws.String(opts.Endpoint),
		Region:       opts.Region,
	}

	return &S3Storage{
		s3:     s3.NewFromConfig(s3Config),
		bucket: opts.Bucket,
		prefix: opts.PathPrefix,
	}, nil
}

func (s *S3Storage) WithTracer(tracer trace.Tracer) TracingStorage {
	return NewTracingStorage(s, tracer)
}

func (s *S3Storage) Purge(ctx context.Context) error {
	if _, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.prefix),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) Get(ctx context.Context, key string) ([]byte, error) {
	filePath := path.Join(s.prefix, key)
	obj, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	value, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := path.Join(s.prefix, key)
	_, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			if awsErr.ErrorCode() == "NotFound" {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]Object, error) {
	searchPath := filepath.Join(s.prefix, prefix)
	objects, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(searchPath),
	})
	if err != nil {
		return nil, err
	}

	results := make([]Object, 0, len(objects.Contents))
	for _, obj := range objects.Contents {
		results = append(results, Object{
			Key:  *obj.Key,
			Size: *obj.Size,
		})
	}

	return results, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	filePath := path.Join(s.prefix, key)
	if _, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) Put(ctx context.Context, key string, data []byte) error {
	filePath := path.Join(s.prefix, key)
	if _, err := s.s3.PutObject(ctx, &s3.PutObjectInput{
		ACL:    types.ObjectCannedACLPublicRead,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return err
	}

	return nil
}
