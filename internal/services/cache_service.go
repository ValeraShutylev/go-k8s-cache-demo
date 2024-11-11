package services

import (
	"cache-demo/internal/errors"
	"cache-demo/utils"
	"log/slog"
	"time"
)

var CLEANUP_DURATION_SEC = 10 * time.Second

type CacheRepositoryInterface interface {
	Get(id string) (any, bool)
	Put(id string, object any, expires uint32) error
	Len() int
	CleanExpiredCacheItems()
	UploadDataToCacheFromStorage(state *utils.GlobalState)
}

type CacheService struct {
	memCache CacheRepositoryInterface
}

func NewCacheService(memCache CacheRepositoryInterface) *CacheService {
	return &CacheService{
		memCache: memCache,
	}
}

func(s *CacheService) GetObjectById(objectId string) (any, error) {
	object, ok := s.memCache.Get(objectId)
	if !ok {
		slog.Error(
			"Object is not found in memory cache",
			slog.String("Key", objectId),
		)
		return nil, errors.ObjectNotFoundException(objectId)
	}
	slog.Info(
		"Object found in memory cache",
		slog.String("Key", objectId),
	)
	return object, nil
}

func(s *CacheService) PutObjectById(objectId string, object any, expires uint32) (isPresent bool) {
	_, ok := s.memCache.Get(objectId)
	if ok {
		isPresent = true
		slog.Info(
			"Object presents in memory cache. Let's update this",
			slog.String("Key", objectId),
		)
	}
	s.memCache.Put(objectId, object, expires)
	slog.Info(
		"Set object to memory cache",
		slog.String("Key", objectId),
	)
	slog.Debug(
		"Found object in memory cache",
		slog.String("Key", objectId),
		slog.Any("Value", object),
		slog.Uint64("Expires at", uint64(expires)),
	)
	return
}

