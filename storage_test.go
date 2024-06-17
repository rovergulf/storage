package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type storageSuite struct {
	suite.Suite

	storage Storage
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(storageSuite))
}

func (s *storageSuite) TestFailBackends() {
	_, fail1 := NewStorage()
	s.Error(fail1, "should return unsupported blockchain error")
	s.Equal(ErrUnsupportedBackends, fail1, "should return unsupported blockchain error")

	_, fail2 := NewStorage(WithBackends("something"))
	s.Error(fail2, "should return unsupported blockchain error")
	s.Equal(ErrUnsupportedBackends, fail2, "should return unsupported blockchain error")

	_, fail3 := NewStorage(WithBackends(S3Backends))
	s.Error(fail3, "should return no s3 config error")
	s.Equal(ErrNoS3ConfigProvided, fail3, "should return no s3 config error")

	_, fail4 := NewStorage(WithBackends(LocalBackends))
	s.Error(fail4, "should return unsafe local storage path errorr")
	s.Equal(ErrLocalStorageUnsafePrefix, fail4, "should return unsafe local storage path error")
}
