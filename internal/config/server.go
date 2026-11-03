package config

import "flag"

var Config struct {
	Address     string
	BaseAddress string
}

func ParseFlags() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "input host")
	flag.StringVar(&Config.BaseAddress, "b", "localhost:8080", "output host")
	flag.Parse()
}
