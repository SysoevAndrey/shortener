package main

import (
	"log"
	"net/http"

	"github.com/SysoevAndrey/shortener/internal/config"
	"github.com/SysoevAndrey/shortener/internal/handler"
	"github.com/SysoevAndrey/shortener/internal/handler/middleware"
	"github.com/SysoevAndrey/shortener/internal/logger"
	"github.com/SysoevAndrey/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	if err := config.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}

	s := repository.NewMemStorage()
	h := handler.NewShortenerHandler(s, config.Config.BaseAddress)
	r := chi.NewRouter()

	r.Use(middleware.WithLogging)
	r.Use(middleware.WithCompression)

	r.Post("/", h.Add)
	r.Post("/api/shorten", h.Shorten)
	r.Get("/{id}", h.Get)

	address := config.Config.Address
	logger.Log.Info("Server started", zap.String("address", address))
	return http.ListenAndServe(address, r)
}
