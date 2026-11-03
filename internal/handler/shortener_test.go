package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SysoevAndrey/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

type storage struct {
	urls map[string]string
}

func newStorage() *storage {
	return &storage{
		urls: make(map[string]string),
	}
}

func (s *storage) Set(url string) (string, error) {
	id := service.GetHash(url)
	if _, exists := s.urls[id]; exists {
		return "", errors.New("url already exists")
	}
	s.urls[id] = url
	return id, nil
}

func (s *storage) Get(hash string) (string, bool) {
	url, found := s.urls[hash]
	return url, found
}

func Test_shortenerHandler_Add(t *testing.T) {
	tests := []struct {
		name    string
		storage *storage
		url     string
		code    int
	}{
		{
			name:    "positive",
			storage: newStorage(),
			url:     "https://ya.ru",
			code:    201,
		},
		{
			name:    "empty body",
			storage: newStorage(),
			code:    400,
		},
		{
			name:    "not a url in body",
			storage: newStorage(),
			url:     "hello, world",
			code:    400,
		},
		{
			name:    "duplicated url",
			storage: &storage{urls: map[string]string{service.GetHash("https://ya.ru"): "https://ya.ru"}},
			url:     "https://ya.ru",
			code:    400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			h := NewShortenerHandler(tt.storage)

			h.Add(w, request)
			res := w.Result()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func Test_shortenerHandler_Get(t *testing.T) {
	type want struct {
		url  string
		code int
	}

	tests := []struct {
		name    string
		storage *storage
		id      string
		want    want
	}{
		{
			name:    "positive",
			storage: &storage{urls: map[string]string{"Qkd1-Ska": "https://ya.ru"}},
			id:      "Qkd1-Ska",
			want:    want{url: "https://ya.ru", code: 307},
		},
		{
			name:    "not found",
			storage: newStorage(),
			id:      "Qkd1-Ska",
			want:    want{code: 404},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.SetPathValue("id", tt.id)

			w := httptest.NewRecorder()
			h := NewShortenerHandler(tt.storage)

			h.Get(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.url, res.Header.Get("Location"))
		})
	}
}
