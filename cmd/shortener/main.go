package main

import (
	"net/http"

	"github.com/SysoevAndrey/shortener/internal/handler"
	"github.com/SysoevAndrey/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := repository.NewMemStorage()
	h := handler.NewShortenerHandler(s)
	r := chi.NewRouter()

	r.Post("/", h.Add)
	r.Get("/{id}", h.Get)

	return http.ListenAndServe(`:8080`, r)
}
