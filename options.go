package storage

import (
	"errors"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

var (
	ErrNoS3ConfigProvided       = errors.New("no S3 config provided")
	ErrUnsupportedBackends      = errors.New("unsupported backends")
	ErrLocalStorageUnsafePrefix = errors.New("path prefix must contain at least two segments")
)

type Options struct {
	backends   Backends
	pathPrefix string
	tracer     trace.Tracer
	s3opts     *S3Options
}

func (o *Options) validateAndFix() error {
	switch o.backends {
	case S3Backends:
		if o.s3opts == nil {
			return ErrNoS3ConfigProvided
		}

		if err := o.s3opts.Validate(); err != nil {
			return err
		}
	case LocalBackends:
		// to prevent uncontrolled file writes on local machines
		if len(strings.Split(o.pathPrefix, "/")) < 2 {
			return ErrLocalStorageUnsafePrefix
		}
	default:
		return ErrUnsupportedBackends
	}

	return nil
}

type Option func(o *Options)

func WithBackends(backends Backends) Option {
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
