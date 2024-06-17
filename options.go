package storage

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"path/filepath"
)

type Options struct {
	backends   StorageType
	pathPrefix string
	tracer     trace.Tracer
	s3opts     *S3Options
}

func (o *Options) validateAndFix() error {
	if ok := supportedStorageTypes[o.backends]; !ok {
		return fmt.Errorf("backends not supported or not specified: '%s'", o.backends)
	}

	if o.backends == LocalBackends {
		// to prevent uncontrolled file writes on local machines
		if len(filepath.SplitList(o.pathPrefix)) < 2 {
			return fmt.Errorf("path prefix must contain at least two segments")
		}
	}

	return nil
}

type Option func(o *Options)

func WithBackends(backends StorageType) Option {
	return func(o *Options) {
		o.backends = backends
	}
}

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.pathPrefix = prefix
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return func(o *Options) {
		o.tracer = tracer
	}
}

func WithS3Config(s3opts *S3Options) Option {
	return func(o *Options) {
		o.s3opts = s3opts
	}
}
