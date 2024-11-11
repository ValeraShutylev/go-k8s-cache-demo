package repository

import (
	"cache-demo/internal/storage"
	"cache-demo/utils"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

const (
	BUFFER_SIZE = 10
)

var CLEANUP_DURATION_SEC = 10 * time.Second

type item struct {
	value   any
	expires uint32
}

type MemoryCache struct {
	items map[string]*item
	mu    *sync.RWMutex
	storage *storage.Storage
}

func NewMemoryCache(storage *storage.Storage) *MemoryCache {
	return &MemoryCache{items: make(map[string]*item), mu: &sync.RWMutex{}, storage: storage}
}

func(mem *MemoryCache) Len() int {
    return len(mem.items)
}

func(item *item) expired(time uint32) bool {
	if item.expires == 0 {
		return true
	}
	return time > item.expires
}

func(mem *MemoryCache) Get(key string) (value any, ok bool) {
	mem.mu.RLock()
    defer mem.mu.RUnlock()
	itm, ok := mem.items[key]
	if !ok {
		slog.Info(
			"Object is not found in memory cache. Try to load from storage",
			slog.String("Key", key),
		)
		var storageData storage.StorageData
		content, persisted := mem.storage.Get(key)
		if !persisted {
			return nil, false
		}
		_ = json.Unmarshal(content, &storageData)
		mem.items[key] = &item{
			value: storageData.Value,
			expires: storageData.Expires,
		}
		slog.Info(
			"Object is not found in storage",
			slog.String("Key", key),
		)
		return storageData.Value, true
	}
    value = itm.value
	slog.Info(
		"Found object in memory cache",
		slog.String("Key", key),
	)
	slog.Debug(
		"Found object in memory cache",
		slog.String("Key", key),
		slog.Any("Value", value),
	)
    return value, true
}

func(mem *MemoryCache) Put(id string, object any, expires uint32) error {
    mem.mu.Lock()
	defer mem.mu.Unlock()

	err := mem.storage.Put(id, object, expires)
	if err != nil {
		return err
	}
    mem.items[id] = &item{
		value:   object,
		expires: expires,
	}
	slog.Info(
		"Successfully set object to memory cache", 
		slog.String("Key", id),
	)
	slog.Debug(
		"Successfully set object to memory cache", 
		slog.String("Key", id),
		slog.Any("Value", object),
		slog.Uint64("Expires at", uint64(expires)),
	)
	return nil
}

func(mem *MemoryCache) clean() {	
 	for k, v := range mem.items {
 		if v.expired(uint32(time.Now().Local().Unix())) {
 			slog.Info("Object has been expired",
				slog.String("Key", k),
				slog.Int64("Expiration time", time.Now().Local().Unix()),
			)
			mem.storage.Delete(k)
 			delete(mem.items, k)
 		}
 	}
}

func(mem *MemoryCache) CleanExpiredCacheItems() {
	t := time.NewTicker(CLEANUP_DURATION_SEC)
	defer t.Stop()
	slog.Info(
		"Starting Expired data cleanup service",
	)
	for range t.C {
		mem.mu.RLock()
		mem.clean()
		mem.mu.RUnlock()
	}
}

func(mem *MemoryCache) UploadDataToCacheFromStorage(state *utils.GlobalState) {
	state.SetState(utils.STARTED)
	files, _ := mem.storage.GetFilenames()

	if len(files) > 0 {
		state.SetState(utils.IN_PROGRESS)
		slog.Info(
			"Read data from storage",
			slog.Int("Files count", len(files)),
		)
		fileChannel := make(chan map[string]*item, BUFFER_SIZE)
		var wg sync.WaitGroup

		for _, file := range files {
			wg.Add(1)

			go func() {
				defer wg.Done()

				content, ok := mem.storage.Get(file.Name())
				if ok {
					var storageData storage.StorageData
					_ = json.Unmarshal(content, &storageData)
					slog.Info(
						"Get data from file",
						slog.String("File", file.Name()),
					)
					slog.Debug(
						"Get data from file",
						slog.String("File", file.Name()),
						slog.Any("Data", storageData),
					)
					fileChannel <- map[string]*item{file.Name() : {
						value: storageData.Value,
						expires: storageData.Expires,
					}}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(fileChannel)
		}()

		for fileData := range fileChannel {
			for key, value := range fileData {
				slog.Info(
					"Set object from storage to memory cache",
					slog.String("Object ID", key),
				)
				slog.Debug(
					"Set object from storage to memory cache",
					slog.String("Object ID", key),
					slog.Any("Object", value),
				)
				mem.items[key] = value
			}
		}
	} else {
		slog.Info(
			"Directory is empty. Data load had skipped",
			slog.String("Directory", storage.CACHE_FILES_PATH),
		)
	}
	state.SetState(utils.DONE)
}	
