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

	_, fail3 := NewStorage(WithBackends("something"))
	s.Error(fail3, "should return unsupported blockchain error")
	s.Equal(ErrUnsupportedBackends, fail3, "should return unsupported blockchain error")

}
