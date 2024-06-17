package storage

import (
	"context"
	"errors"
)

const (
	S3Backends    = "s3"
	LocalBackends = "file"
	//GCSBackends   = "gcs" // To be implemented
)

var ErrNotExists = errors.New("file not exists")

// Storage represents main interface for file storage as key/value database
type Storage interface {
	Put(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]Object, error)
	Purge(ctx context.Context) error
}

type Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	Data []byte `json:"data"` // not currently used in either list queries or the for get method
}
