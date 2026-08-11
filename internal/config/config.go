package config

import "time"

type Config struct {
	Port string
	DBPath          string
	SessionDuration time.Duration
}

func Load() Config {
	return Config{
		Port: ":8081",
		DBPath:          "./realTime.db",
		SessionDuration: 24 * time.Hour,
	}
}
