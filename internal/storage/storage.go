package storage

import (
	"io"
	"os"
)

var basePath string = "data/images"

type Storage interface {
	Save(data []byte) (string, error)
	Get(path string) ([]byte, error)
}

type ImageStorage struct {
	BasePath string
}

func NewImageStorage() *ImageStorage {
	return &ImageStorage{
		BasePath: basePath,
	}
}

func (is *ImageStorage) Save(data []byte, name string) (string, error) {
	path := is.BasePath + "/" + name
	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (is *ImageStorage) Get(path string) ([]byte, error) {
	file, err := os.Open(is.BasePath + "/" + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
