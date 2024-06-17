package storage

import (
	"context"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/trace"
)

type TracingStorage interface {
	Storage
}

type tracedStorage struct {
	storage Storage

	tracer trace.Tracer
}

func NewTracingStorage(storage Storage, tracer trace.Tracer) TracingStorage {
	return &tracedStorage{
		storage: storage,
		tracer:  tracer,
	}
}

func (s *tracedStorage) Get(ctx context.Context, key string) ([]byte, error) {
	if s.tracer != nil {
		var span trace.Span
		ctx, span = s.tracer.Start(ctx, "get",
			trace.WithAttributes(
				attribute.String("key", key),
			),
		)
		defer span.End()
	}

	return s.storage.Get(ctx, key)
}

//func (s *tracedStorage) GetMultiple(ctx context.Context, filePaths []string) ([]Object, error) {
//	return s.storage.GetMultiple(ctx, filePaths)
//}

func (s *tracedStorage) Exists(ctx context.Context, key string) (bool, error) {
	if s.tracer != nil {
		var span trace.Span
		ctx, span = s.tracer.Start(ctx, "exists",
			trace.WithAttributes(
				attribute.String("key", key),
			),
		)
		defer span.End()
	}

	return s.storage.Exists(ctx, key)
}

func (s *tracedStorage) List(ctx context.Context, prefix string) ([]Object, error) {
	if s.tracer != nil {
		var span trace.Span
		ctx, span = s.tracer.Start(ctx, "list",
			trace.WithAttributes(
				attribute.String("prefix", prefix),
			),
		)
		defer span.End()
	}

	return s.storage.List(ctx, prefix)
}

func (s *tracedStorage) Delete(ctx context.Context, key string) error {
	if s.tracer != nil {
		var span trace.Span
		ctx, span = s.tracer.Start(ctx, "delete",
			trace.WithAttributes(
				attribute.String("key", key),
			),
		)
		defer span.End()
	}

	return s.storage.Delete(ctx, key)
}

func (s *tracedStorage) Put(ctx context.Context, key string, data []byte) error {
	if s.tracer != nil {
		var span trace.Span
		ctx, span = s.tracer.Start(ctx, "put",
			trace.WithAttributes(
				attribute.String("key", key),
			),
		)
		defer span.End()
	}

	return s.storage.Put(ctx, key, data)
}

func (s *tracedStorage) Purge(ctx context.Context) error {
	return s.storage.Purge(ctx)
}
