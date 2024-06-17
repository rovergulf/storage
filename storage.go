package storage

import (
	"context"
	"errors"
	"fmt"
)

type StorageType string

func (s StorageType) String() string {
	return string(s)
}

const (
	S3Backends    StorageType = "s3"
	LocalBackends             = "file"
	//GCSBackends   = "gcs" // To be implemented
)

var supportedStorageTypes = map[StorageType]bool{
	LocalBackends: true,
	S3Backends:    true,
	//GCSBackends:   true,
}

var ErrNotExists = errors.New("file not exists")
var ErrUnsupportedBackends = errors.New("unsupported backends")

// Storage represents main interface for file storage as key/value database
type Storage interface {
	Put(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetMultiple(ctx context.Context, keys []string) ([]Object, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]Object, error)
	Purge(ctx context.Context) error
}

type Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	Data []byte `json:"data"`
}

func NewStorage(options ...Option) (Storage, error) {
	opts := new(Options)
	for _, opt := range options {
		opt(opts)
	}

	err := opts.validateAndFix()
	if err != nil {
		return nil, err
	}

	var s Storage
	switch opts.backends {
	case S3Backends:
		if opts.s3opts == nil {
			return nil, fmt.Errorf("no S3 config provided")
		}

		if err := opts.s3opts.Validate(); err != nil {
			return nil, err
		}

		s, err = NewS3Storage(opts.s3opts)
		if err != nil {
			return nil, err
		}
	case LocalBackends:
		s = NewFileStorage(opts.pathPrefix)
	default:
		return nil, ErrUnsupportedBackends
	}

	if opts.tracer != nil {
		return NewTracingStorage(s, opts.tracer), nil
	}

	return s, nil
}
