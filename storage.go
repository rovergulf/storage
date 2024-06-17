package storage

import (
	"context"
	"errors"
)

type Backends string

func (s Backends) String() string {
	return string(s)
}

const (
	LocalBackends Backends = "file"
	S3Backends    Backends = "s3"
	//GCSBackends   Backends = "gcs" // To be implemented
)

var supportedBackendss = map[Backends]bool{
	LocalBackends: true,
	S3Backends:    true,
	//GCSBackends:   true,
}

var ErrNotExists = errors.New("file not exists")

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
