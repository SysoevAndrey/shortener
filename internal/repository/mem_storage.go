package repository

import (
	"errors"

	"github.com/SysoevAndrey/shortener/internal/service"
)

type memStorage struct {
	urls map[string]string
}

func NewMemStorage() *memStorage {
	return &memStorage{
		urls: make(map[string]string),
	}
}

func (s *memStorage) Set(url string) (string, error) {
	h := service.GetHash(url)
	if _, exists := s.urls[h]; exists {
		return "", errors.New("url already exists")
	}
	s.urls[h] = url
	return h, nil
}

func (s *memStorage) Get(hash string) (string, bool) {
	url, found := s.urls[hash]
	return url, found
}
