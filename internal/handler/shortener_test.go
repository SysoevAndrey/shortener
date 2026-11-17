package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SysoevAndrey/shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storage struct {
	urls map[string]string
}

func newStorage() *storage {
	return &storage{
		urls: make(map[string]string),
	}
}

const (
	host = "http://localhost:8080"
)

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
			r := chi.NewRouter()
			h := NewShortenerHandler(tt.storage, host)

			r.Post("/", h.Add)
			r.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func Test_shortenerHandler_Shorten(t *testing.T) {
	type want struct {
		json string
		code int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive",
			body: `{"url": "https://ya.ru"}`,
			want: want{json: `{"result": "http://localhost:8080/fpCk-cML"}`, code: 201},
		},
		{
			name: "returns 400 on invalid url",
			body: `hello, world`,
			want: want{code: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewReader([]byte(tt.body))
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bodyReader)
			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r := chi.NewRouter()
			h := NewShortenerHandler(newStorage(), host)

			r.Post("/api/shorten", h.Shorten)
			r.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)
			if len(tt.want.json) > 0 {
				assert.JSONEq(t, tt.want.json, string(body))
			}
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
			storage: &storage{urls: map[string]string{"fpCk-cML": "https://ya.ru"}},
			id:      "fpCk-cML",
			want:    want{url: "https://ya.ru", code: 307},
		},
		{
			name:    "not found",
			storage: newStorage(),
			id:      "fpCk-cML",
			want:    want{code: 404},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)

			w := httptest.NewRecorder()
			r := chi.NewRouter()
			h := NewShortenerHandler(tt.storage, host)

			r.Get("/{id}", h.Get)
			r.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.url, res.Header.Get("Location"))
		})
	}
}

func Test_shortenerHandler_processURL(t *testing.T) {
	tests := []struct {
		name    string
		urlRaw  string
		want    string
		wantErr bool
	}{
		{
			name:    "stores url",
			urlRaw:  "https://ya.ru",
			want:    "fpCk-cML",
			wantErr: false,
		},
		{
			name:    "throws error on invalid url",
			urlRaw:  "hello, world",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewShortenerHandler(newStorage(), host)
			hash, err := h.processURL(tt.urlRaw)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.want, hash)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, hash)

			url, found := h.storage.Get(hash)
			require.True(t, found)
			assert.Equal(t, tt.urlRaw, url)
		})
	}
}
