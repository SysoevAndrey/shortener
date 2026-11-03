package config

import (
	"flag"
	"strings"
)

var Config struct {
	Address     string
	BaseAddress string
}

func ParseFlags() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "input host")
	flag.StringVar(&Config.BaseAddress, "b", "localhost:8080", "output host")
	flag.Parse()

	Config.BaseAddress = normalizeHost(Config.BaseAddress)
}

func normalizeHost(host string) string {
	if strings.Contains(host, "://") {
		return host
	}
	return "http://" + host
}
