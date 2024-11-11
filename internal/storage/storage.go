package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)
const (
	CACHE_FILES_PATH = "cachefiles"
	BUFFER_SIZE = 10
)

type Storage struct{}

type StorageData struct {
	Value any
	Expires uint32
}

func init() {
	err := os.MkdirAll(CACHE_FILES_PATH, 0755)
	if err != nil {
		slog.Error(
			"Error during the cache directory creation",
			slog.String("Error", err.Error()),
		)
	}
	slog.Info(
		"Direcory successfully created",
		slog.String("Directory", CACHE_FILES_PATH),
	)
}

func NewStorage() *Storage {
	return &Storage{}
}

func(storage *Storage) Get(key string) (data []byte, ok bool) {
	filePath := filepath.Join(CACHE_FILES_PATH, key)
	data, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		slog.Info("Object not found from storage",
			slog.String("Filename", filePath ),
		)
        return nil, false
    }
	slog.Info(
		"Found object in storage",
		slog.String("File", filePath),
	)
	return data, true
}

func(storage *Storage) Put(key string, value any, expires uint32) error {
	filePath := filepath.Join(CACHE_FILES_PATH, key)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666) 
	if err != nil {
		slog.Error(
			"Cannot open file for writing",
			slog.String("File", filePath),
			slog.String("Error", err.Error()),
		)
		return err
	} 
	defer file.Close()
	
	byteData, _ := json.Marshal(&StorageData{Value: value, Expires: expires})
	_, err = file.Write(byteData)
    if err != nil {
		slog.Error(
			"Error during file writing",
			slog.String("File", filePath),
			slog.String("Error", err.Error()),
		)
        return err
    }
	slog.Info(
		"Successfully put object to storage",
		slog.String("File", filePath),
	)
	return nil
}

func(storage *Storage) Delete(key string) error {
	filePath := filepath.Join(CACHE_FILES_PATH, key)
	err := os.Remove(filePath)
	if err != nil {
		slog.Error(
			"Error during file removal",
			slog.String("File", filePath),
			slog.String("Error", err.Error()),
		)
		return err
	}
	slog.Info(
		"File successfully removed",
		slog.String("File", filePath),
	)
	return nil
}

func(storage *Storage) GetFilenames() ([]fs.DirEntry, error) {
	files, err := os.ReadDir(CACHE_FILES_PATH)
	slog.Info(
		"Read all filenames from directory",
		slog.String("Directory", CACHE_FILES_PATH),
	)
	if err != nil {
		slog.Error(
			"Error during directory reading",
			slog.String("Error", err.Error()),
		)
		return nil, err
	}
	return files, nil
}

