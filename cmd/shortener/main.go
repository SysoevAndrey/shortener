package main

import (
	"fmt"
	"net/http"

	"github.com/SysoevAndrey/shortener/internal/config"
	"github.com/SysoevAndrey/shortener/internal/handler"
	"github.com/SysoevAndrey/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	s := repository.NewMemStorage()
	h := handler.NewShortenerHandler(s, config.Config.BaseAddress)
	r := chi.NewRouter()

	r.Post("/", h.Add)
	r.Get("/{id}", h.Get)

	address := config.Config.Address
	fmt.Println("Server started on:", address)
	return http.ListenAndServe(address, r)
}
