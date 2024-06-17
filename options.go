package storage

import "go.opentelemetry.io/otel/trace"

type Options struct {
	pathPrefix string
	tracer     trace.Tracer
	s3opts     *S3Options
}

type Option func(o *Options)

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
