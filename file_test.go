package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type fileStorageSuite struct {
	suite.Suite

	fs FileStorage
}

func (s *fileStorageSuite) SetupSuite() {
	s.fs = FileStorage{
		dir: "./tmp-test",
	}
}

func (s *fileStorageSuite) TearDownSuite() {
	os.RemoveAll(s.fs.dir)
}

func TestFileStorage(t *testing.T) {
	suite.Run(t, new(fileStorageSuite))
}

func (s *fileStorageSuite) TestFileStorage() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testData1 := `{"example": "hello world 1"}`
	testData2 := `{"example": "hello world 2"}`

	file1Name := "example1.json"
	file2Name := "example2.json"

	err1 := s.fs.Put(ctx, file1Name, []byte(testData1))
	s.NoError(err1, "shouldn't return any error on put file 1")
	fileInfo1, errStat1 := os.Stat(filepath.Join(s.fs.dir, file1Name))
	s.NoError(errStat1, "shouldn't return any error on stat file 1")
	s.NotNil(fileInfo1, "should return file 1 info")

	err2 := s.fs.Put(ctx, file2Name, []byte(testData2))
	s.NoError(err2, "shouldn't return any error on put file 2")
	fileInfo2, errStat2 := os.Stat(filepath.Join(s.fs.dir, file2Name))
	s.NoError(errStat2, "shouldn't return any error on stat file 2")
	s.NotNil(fileInfo2, "should return file 2 info")

	exists, err3 := s.fs.Exists(ctx, file2Name)
	s.NoError(err3, "shouldn't return any error on exists check")
	s.Equal(true, exists, "should confirm file existence")

	files, listErr := s.fs.List(ctx, s.fs.dir)
	s.NoError(listErr, "shouldn't return any error on files list")
	s.Greater(len(files), 0, "should return at least one record")

	deleteErr := s.fs.Delete(ctx, file1Name)
	s.NoError(deleteErr, "shouldn't return any error on delete file1Name")
	_, fileStatErr := os.Stat(file1Name)
	s.Error(os.ErrNotExist, fileStatErr, "should return os.ErrNotExists after delete")
}

func (s *fileStorageSuite) TestGetFile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testData3 := `{"example": "hello world 3"}`

	file3Name := "example3.json"

	err4 := s.fs.Put(ctx, file3Name, []byte(testData3))
	s.NoError(err4, "shouldn't return any error on put file 3")

	payload, errGet := s.fs.Get(ctx, file3Name)
	s.NoError(errGet, "shouldn't return any error on get file 3")
	s.GreaterOrEqual(len(payload), len(testData3), "payload length should be greater than 0")
}
