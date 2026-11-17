package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v6"
)

var Config struct {
	Address     string `env:"SERVER_ADDRESS"`
	BaseAddress string `env:"BASE_URL"`
}

func ParseFlags() error {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "input host")
	flag.StringVar(&Config.BaseAddress, "b", "localhost:8080", "output host")
	flag.Parse()

	if err := env.Parse(&Config); err != nil {
		return err
	}

	Config.BaseAddress = normalizeHost(Config.BaseAddress)

	return nil
}

func normalizeHost(host string) string {
	if strings.Contains(host, "://") {
		return host
	}
	return "http://" + host
}
