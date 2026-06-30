package config

import "time"

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
	Version         string
}

func Default() Config {
	return Config{
		Addr:            ":8080",
		ShutdownTimeout: 30 * time.Second,
		Version:         "dev",
	}
}

func Load() Config {
	return Default()
}
