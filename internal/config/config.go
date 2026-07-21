package config

import "time"

type Config struct {
	PortMux1        string
	PortMux2        string
	DBPath          string
	SessionDuration time.Duration
}

func Load() Config {
	return Config{
		PortMux1:        ":8082",
		PortMux2:        ":8081",
		DBPath:          "./realTime.db",
		SessionDuration: 24 * time.Hour,
	}
}
