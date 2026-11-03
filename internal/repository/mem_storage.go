package repository

import (
	"errors"
	"sync"

	"github.com/SysoevAndrey/shortener/internal/service"
)

type memStorage struct {
	urls map[string]string
	mu   sync.Mutex
}

func NewMemStorage() *memStorage {
	return &memStorage{
		urls: make(map[string]string),
	}
}

func (s *memStorage) Set(url string) (string, error) {
	h := service.GetHash(url)
	s.mu.Lock()
	if _, exists := s.urls[h]; exists {
		return "", errors.New("url already exists")
	}
	s.urls[h] = url
	s.mu.Unlock()
	return h, nil
}

func (s *memStorage) Get(hash string) (string, bool) {
	s.mu.Lock()
	url, found := s.urls[hash]
	s.mu.Unlock()
	return url, found
}
