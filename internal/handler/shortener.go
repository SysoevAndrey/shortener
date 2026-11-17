package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/SysoevAndrey/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

type shortenRequestBody struct {
	URL string `json:"url"`
}

type shortenResponseBody struct {
	Result string `json:"result"`
}

func (h *shortenerHandler) Add(w http.ResponseWriter, r *http.Request) {
	bodyRaw, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlString := string(bodyRaw)
	hash, err := h.processURL(urlString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.host + "/" + hash))
}

func (h *shortenerHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var requestPayload shortenRequestBody

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&requestPayload); err != nil {
		logger.Log.Debug("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := h.processURL(requestPayload.URL)
	if err != nil {
		logger.Log.Debug("error processing url", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responsePayload := shortenResponseBody{Result: h.host + "/" + hash}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(responsePayload); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
	logger.Log.Debug("sending HTTP 201 response")
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

func (h *shortenerHandler) processURL(urlRaw string) (string, error) {
	_, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		return "", err
	}

	hash, err := h.storage.Set(urlRaw)
	if err != nil {
		return "", err
	}

	return hash, nil
}
