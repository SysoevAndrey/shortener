package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type shortenerHandler struct {
	storage URLStorage
	host    string
}

func NewShortenerHandler(storage URLStorage, host string) *shortenerHandler {
	return &shortenerHandler{
		storage: storage,
		host:    host,
	}
}

type URLStorage interface {
	Set(url string) (string, error)
	Get(hash string) (string, bool)
}

func (h *shortenerHandler) Add(w http.ResponseWriter, r *http.Request) {
	bodyRaw, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlString := string(bodyRaw)
	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := h.storage.Set(urlString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.host + "/" + hash))
}

func (h *shortenerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok := h.storage.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
