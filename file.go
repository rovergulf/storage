package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/trace"
)

var (
	defaultDirStoragePath = filepath.Join(os.TempDir(), "go-storage")
)

type FileStorage struct {
	dir string

	tracer trace.Tracer
}

func NewFileStorage(dir string) *FileStorage {
	return &FileStorage{dir: dir}
}

func (s *FileStorage) WithTracer(tracer trace.Tracer) TracingStorage {
	return NewTracingStorage(s, tracer)
}

func (s *FileStorage) Purge(ctx context.Context) error {
	return os.RemoveAll(s.dir)
}

func (s *FileStorage) Get(ctx context.Context, key string) ([]byte, error) {
	filePath := filepath.Join(s.dir, key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *FileStorage) GetMultiple(ctx context.Context, filePaths []string) ([]Object, error) {
	var res []Object
	for _, filePath := range filePaths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		res = append(res, Object{
			Key:  filePath,
			Data: data,
		})
	}

	return res, nil
}

func (s *FileStorage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(s.dir, key)
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func (s *FileStorage) List(ctx context.Context, prefix string) ([]Object, error) {
	dirData, err := os.ReadDir(prefix)
	if err != nil {
		return nil, err
	}

	results := make([]Object, 0, len(dirData))
	for _, dirEntry := range dirData {
		results = append(results, Object{Key: dirEntry.Name()})
	}

	return results, nil
}

func (s *FileStorage) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(s.dir, key)
	return os.Remove(filePath)
}

func (s *FileStorage) Put(ctx context.Context, key string, data []byte) error {
	if _, err := os.Stat(s.dir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(s.dir, os.ModePerm); err != nil {
			return err
		}
	}

	keyDir := filepath.Join(s.dir, filepath.Dir(key))
	if len(keyDir) > 0 {
		if _, err := os.Stat(keyDir); errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(keyDir, os.ModePerm); err != nil {
				return err
			}
		}
	}

	filePath := filepath.Join(s.dir, key)
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}
