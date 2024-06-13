package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type s3Suite struct {
	suite.Suite

	opts *S3Options

	s3 *S3Storage
}

func (s *s3Suite) SetupSuite() {
	accessKey := os.Getenv("TEST_S3_ACCESS_KEY")
	secretKey := os.Getenv("TEST_S3_SECRET_KEY")
	endpoint := os.Getenv("TEST_S3_ENDPOINT")
	region := os.Getenv("TEST_S3_REGION")
	bucket := os.Getenv("TEST_S3_BUCKET")

	s.opts = &S3Options{
		Endpoint:   endpoint,
		Region:     region,
		Key:        accessKey,
		Secret:     secretKey,
		Bucket:     bucket,
		PathPrefix: "/tests/s3-storage",
	}

	if err := s.opts.Validate(); err != nil {
		s.T().Skip()
		return
	}

	s3, err := NewS3Storage(s.opts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.s3 = s3
}

func (s *s3Suite) TearDownSuite() {
	if s.s3 != nil {
		_ = s.s3.Purge(context.Background())
	}
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(s3Suite))
}

func (s *s3Suite) TestS3Suite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testData1 := `{"example": "hello world 1"}`
	testData2 := `{"example": "hello world 2"}`

	file1Name := "example1.json"
	file2Name := "example2.json"

	err1 := s.s3.Put(ctx, file1Name, []byte(testData1))
	s.NoError(err1, "shouldn't return any error on put file 1")
	if err1 != nil {
		return
	}
	exists1, errStat1 := s.s3.Exists(ctx, file1Name)
	s.NoError(errStat1, "shouldn't return any error on stat file 1")
	s.Equal(true, exists1, "should exist")
	if errStat1 != nil {
		return
	}

	err2 := s.s3.Put(ctx, file2Name, []byte(testData2))
	s.NoError(err2, "shouldn't return any error on put file 2")
	if err2 != nil {
		return
	}

	exists, err3 := s.s3.Exists(ctx, file2Name)
	s.NoError(err3, "shouldn't return any error on exists check")
	s.Equal(true, exists, "should confirm file existence")
	if err3 != nil {
		return
	}

	files, listErr := s.s3.List(ctx, "")
	s.NoError(listErr, "shouldn't return any error on files list")
	s.Greater(len(files), 0, "should return at least one record")
	if listErr != nil {
		return
	}

	deleteErr := s.s3.Delete(ctx, file1Name)
	s.NoError(deleteErr, "shouldn't return any error on delete file1Name")
}

func (s *s3Suite) TestNotExists() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	unexistingFile := "/abra/cadabra/non-existent.txt"
	exists, err := s.s3.Exists(ctx, unexistingFile)
	s.NoError(err, "shouldn't return any error on exists check")
	s.Equal(false, exists, "should not exist")
}

func (s *s3Suite) TestGetFile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testData3 := `{"example": "hello world 3"}`

	file3Name := "example3.json"

	err4 := s.s3.Put(ctx, file3Name, []byte(testData3))
	s.NoError(err4, "shouldn't return any error on put file 3")

	payload, errGet := s.s3.Get(ctx, file3Name)
	s.NoError(errGet, "shouldn't return any error on get file 3")
	s.GreaterOrEqual(len(payload), len(testData3), "payload length should be greater than 0")
}
