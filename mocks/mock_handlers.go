package mocks

import (
	"github.com/stretchr/testify/mock"
)

type CacheServiceMock struct {
	mock.Mock
}

func NewServiceMock() *CacheServiceMock {
	return &CacheServiceMock{}
}

func (s *CacheServiceMock) GetObjectById(objectId string) (any, error) {
	args := s.Called(objectId)
	return args.Get(0), args.Error(1)
}

func (s *CacheServiceMock) PutObjectById(objectId string, object any, expired uint32) bool {
	args := s.Called(objectId, object, expired)
	return args.Bool(0)
}