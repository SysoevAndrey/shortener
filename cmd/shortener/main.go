package main

import (
	"net/http"

	"github.com/SysoevAndrey/shortener/internal/handler"
	"github.com/SysoevAndrey/shortener/internal/repository"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := repository.NewMemStorage()
	h := handler.NewShortenerHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Add)
	mux.HandleFunc("GET /{id}", h.Get)

	return http.ListenAndServe(`:8080`, mux)
}
